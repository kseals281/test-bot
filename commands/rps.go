package commands

import (
  // "fmt"
  // "os"
  // s "strings"
  // "sync"

  "github.com/bwmarrin/discordgo"
)

func RPSHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
  // Webhook messages do not contain a full author (i.e. non-user author) and
  // should be ignored as a result
  if message.WebhookID != "" {
    return
  }

  if message.MentionEveryone {
    discord.ChannelMessageSend(message.ChannelID,
      "You may not challenge everyone to a game of rock, paper, scissors")
    return
  } else if len(message.MentionRoles) > 0 {
    discord.ChannelMessageSend(message.ChannelID,
      "You may not challenge a role to a game of rock, paper, scissors")
    return
  } else if len(message.Mentions) > 1 {
    discord.ChannelMessageSend(message.ChannelID,
      "You can only challenge one person at a time to rock, paper, scissors")
    return
  } else if len(message.Mentions) == 0 {
    discord.ChannelMessageSend(message.ChannelID,
      "You must select at least one opponent for rock, paper, scissors")
    return
  }
  // TODO: Add check for user challenging themselves

  dmChannels := make(chan []*discordgo.Channel)
  // TODO: Add timeout for player choices

  go rpsContactPlayers(discord, message, dmChannels)

  // discord.AddHandlerOnce(rpsGetReaction)

}

func rpsContactPlayers(discord *discordgo.Session, message *discordgo.MessageCreate, dmChannels chan []*discordgo.Channel) {
  player_2 := message.Mentions[0]
  player_1 := message.Author

  dmP1, err := discord.UserChannelCreate(player_1.ID)
  if err != nil {
    errCheck("Unable to create direct message with the challenger", err)
    return
  }
  dmP1.Name = player_1.Username + "1"

  dmP2, err := discord.UserChannelCreate(player_2.ID)
  if err != nil {
    errCheck("Unable to create direct message with the opponent", err)
    return
  }
  dmP2.Name = player_2.Username + "2"

  p1Message := "__**You have challenged " + player_2.Username + " to a game of rock paper scissors!!!**__\n"
  p2Message := "__**" + player_1.Username + " has challenged you to a match of rock paper scissors!!!**__\n"

  commonMessage := (":fist::raised_hand::v:\t:fist::raised_hand::v:\t:fist::raised_hand::v:\t" +
              ":fist::raised_hand::v:\t:fist::raised_hand::v:" +
              "\n\n" +
              "Rules: You must react to this message with either the rock (:fist:)," +
              " paper (:raised_hand:), or scissors (:v:) emoji.\n" +
              "- Rock beats scissors\n" +
              "- Scissors beats paper\n" +
              "- Paper beats rock")

  discord.ChannelMessageSend(dmP1.ID, (p1Message + commonMessage))
  discord.ChannelMessageSend(dmP2.ID, (p2Message + commonMessage))

  dm := make([]*discordgo.Channel, 2)
  dm[0] = dmP1
  dm[1] = dmP2
  dmChannels <- dm

}

// func rpsGetReaction(discord *discordgo.Session, dmChannels *discordgo.Channel) {
//   <- dmChannels
// }