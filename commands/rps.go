package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sync"
	"time"
)

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
	// TODO: Add timeout for player choices

	p1, p2 := rpsContactPlayers(discord, message)
	if p2 == nil {
		errCheck("Error Contacting players", nil)
	}

	var wg sync.WaitGroup
	reactions := make(chan *discordgo.MessageReaction)
	wg.Add(2)
	go rpsWaitForReaction(&wg, discord, p1, reactions)
	go rpsWaitForReaction(&wg, discord, p2, reactions)

	wg.Wait()
	fmt.Println("Finished waiting")
	var roshambo []*discordgo.MessageReaction
	for i := 0; i < 2; i++ {
		roshambo[i] = <-reactions
	}
	fmt.Println(roshambo)
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

func rpsWaitForReaction(wg *sync.WaitGroup, s *discordgo.Session, message *discordgo.Message, reactions chan *discordgo.MessageReaction) {
	// TODO: Do this in a more efficient manner
	if len(message.Reactions) > 0 {
		_ = s.MessageReactionsRemoveAll(message.ChannelID, message.ID)
	}
	limiter := time.Tick(1000 * time.Millisecond)
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
				fmt.Println("GOT A REACTION")
				mr := messageReaction(message.Reactions[0], message)
				reactions <- mr
				wg.Done()
				return
			}
		}
		<-limiter
	}
}

func messageReaction(reaction *discordgo.MessageReactions, message *discordgo.Message) *discordgo.MessageReaction {
	mr := discordgo.MessageReaction{UserID: message.Author.ID, MessageID: message.ID, Emoji: *reaction.Emoji,
		ChannelID: message.ChannelID, GuildID: message.GuildID}
	return &mr
}
