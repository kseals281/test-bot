package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type roshambo struct {
	choice *discordgo.MessageReaction
	user   *discordgo.User
}

func RPSHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Webhook messages do not contain a full author (i.e. non-user author) and
	// should be ignored as a result
	if message.WebhookID != "" {
		return
	}

	var err error
	var msg string

	if message.MentionEveryone {
		_, err = discord.ChannelMessageSend(message.ChannelID,
			"You may not challenge everyone to a game of rock, paper, scissors")
		msg = "Unable to sent cannot challenge everyone message"
	} else if len(message.MentionRoles) > 0 {
		_, err = discord.ChannelMessageSend(message.ChannelID,
			"You may not challenge a role to a game of rock, paper, scissors")
		msg = "Error sending role challenge denial"
	} else if len(message.Mentions) > 1 {
		_, err = discord.ChannelMessageSend(message.ChannelID,
			"You can only challenge one person at a time to rock, paper, scissors")
		msg = "Unable to send only one opponent message"
	} else if len(message.Mentions) == 0 {
		_, err = discord.ChannelMessageSend(message.ChannelID,
			"You must select at least one opponent for rock, paper, scissors")
		msg = "Unable to send at least one opponent error"
	} else if message.Mentions[0].ID == message.Author.ID {
		_, err = discord.ChannelMessageSend(message.ChannelID,
			"You cannot challenge yourself to a game of rock paper scissors")
		msg = "Unable to send no challenging yourself message"
	}

	if err != nil {
		errCheck(msg, err)
		return
	}
	if len(msg) > 0 {
		return
	}

	// TODO: Add timeout for player choices

	p1, p2 := rpsContactPlayers(discord, message)
	if p2 == nil {
		errCheck("Error Contacting players", nil)
	}

	reactions := make(chan *roshambo)
	go rpsWaitForReaction(discord, p1, reactions)
	go rpsWaitForReaction(discord, p2, reactions)

	var rpsChoices [2]*roshambo
	for i := 0; i < 2; i++ {
		rpsChoices[i] = <-reactions
	}
	close(reactions)

	rpsWinner(discord, message, rpsChoices)
}

func rpsContactPlayers(discord *discordgo.Session, message *discordgo.MessageCreate) (*discordgo.Message, *discordgo.Message) {
	player2 := message.Mentions[0]
	player1 := message.Author

	dmP1, err := discord.UserChannelCreate(player1.ID)
	if err != nil {
		errCheck("Unable to create direct message with the challenger", err)
		return nil, nil
	}
	dmP1.Name = player1.Username + "1"

	dmP2, err := discord.UserChannelCreate(player2.ID)
	if err != nil {
		errCheck("Unable to create direct message with the opponent", err)
		return nil, nil
	}
	dmP2.Name = player2.Username + "2"

	p1Message := "__**You have challenged " + player2.Username + " to a game of rock paper scissors!!!**__\n"
	p2Message := "__**" + player1.Username + " has challenged you to a match of rock paper scissors!!!**__\n"

	commonMessage := ":fist::raised_hand::v:\t:fist::raised_hand::v:\t:fist::raised_hand::v:\t" +
		":fist::raised_hand::v:\t:fist::raised_hand::v:" +
		"\n\n" +
		"Rules: You must react to this message with either the rock (:fist:)," +
		" paper (:raised_hand:), or scissors (:v:) emoji.\n" +
		"- Rock beats scissors\n" +
		"- Scissors beats paper\n" +
		"- Paper beats rock"

	p1, err := discord.ChannelMessageSend(dmP1.ID, p1Message+commonMessage)
	if err != nil {
		errCheck("Error sending message to challenger", err)
		return nil, nil
	}
	p2, err := discord.ChannelMessageSend(dmP2.ID, p2Message+commonMessage)
	if err != nil {
		errCheck("Error sending message to opponent", err)
		return nil, nil
	}
	return p1, p2
}

func rpsWaitForReaction(s *discordgo.Session, message *discordgo.Message, reactions chan *roshambo) {
	// TODO: Rewrite entire function using a test driven approach
	// TODO: Do this using Emoji.APIName
	if len(message.Reactions) > 0 {
		_ = s.MessageReactionsRemoveAll(message.ChannelID, message.ID)
	}
	limiter := time.Tick(1000 * time.Millisecond)
	// For some reason the ✌️ emoji is not recognized unless backspace is pressed once
	choices := [3]string{"✊", "✋", "✌"}
	for i := 0; i < 500; i++ {
		message, _ = s.ChannelMessage(message.ChannelID, message.ID) // Refresh message
		if len(message.Reactions) == 0 {
			continue
		}
		// TODO: Remove all non rps reactions
		fmt.Printf("Tick %d\tReactions: %v\n", i, message.Reactions[0].Emoji)
		for _, emoji := range choices {
			users, err := s.MessageReactions(message.ChannelID, message.ID, emoji, 1)
			errCheck("Error getting number of users who reacted", err)
			if len(users) > 0 {
				mr := messageReaction(message.Reactions[0], message)
				r := new(roshambo)
				r.user = users[0]
				r.choice = mr
				reactions <- r
				return
			}
		}
		<-limiter
	}
}

func rpsResults(rpsChoices [2]*roshambo) string {
	if rpsChoices[0].choice.Emoji.Name == "✊" {
		if rpsChoices[1].choice.Emoji.Name == "✌" {
			return rpsChoices[0].user.ID
		} else if rpsChoices[1].choice.Emoji.Name == "✋" {
			return rpsChoices[1].user.ID
		}
	} else if rpsChoices[0].choice.Emoji.Name == "✋" {
		if rpsChoices[1].choice.Emoji.Name == "✊" {
			return rpsChoices[0].user.ID
		} else if rpsChoices[1].choice.Emoji.Name == "✌" {
			return rpsChoices[1].user.ID
		}
	} else if rpsChoices[0].choice.Emoji.Name == "✌" {
		if rpsChoices[1].choice.Emoji.Name == "✋" {
			return rpsChoices[0].user.ID
		} else if rpsChoices[1].choice.Emoji.Name == "✊" {
			return rpsChoices[1].user.ID
		}
	}
	return ""
}

func rpsWinner(discord *discordgo.Session, message *discordgo.MessageCreate, rpsChoices [2]*roshambo) {
	winnerID := rpsResults(rpsChoices)
	var players []*discordgo.User
	for _, choice := range rpsChoices {
		players = append(players, choice.user)
	}

	var winner, loser *discordgo.User
	if winnerID == players[0].ID {
		winner = players[0]
		loser = players[1]
	} else {
		winner = players[1]
		loser = players[0]
	}

	var r string
	if winnerID == "" {
		r = fmt.Sprintf("> There is no winner between %s and %s\n", players[0].Username, players[1].Username)
	} else {
		r = fmt.Sprintf("> **%s has won the rock paper scissors match against %s**\n", winner.Username, loser.Username)
	}
	p := fmt.Sprintf("> %s chose %s\n", rpsChoices[0].user.Username, rpsChoices[0].choice.Emoji.MessageFormat())
	s := fmt.Sprintf("> %s chose %s\n", rpsChoices[1].user.Username, rpsChoices[1].choice.Emoji.MessageFormat())

	_, err := discord.ChannelMessageSend(message.ChannelID, r+p+s)
	errCheck("Failed sending winner of rps game to chat", err)
}

func messageReaction(reaction *discordgo.MessageReactions, message *discordgo.Message) *discordgo.MessageReaction {
	mr := discordgo.MessageReaction{UserID: message.Author.ID, MessageID: message.ID, Emoji: *reaction.Emoji,
		ChannelID: message.ChannelID, GuildID: message.GuildID}
	return &mr
}
