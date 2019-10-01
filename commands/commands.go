package commands

import (
  "fmt"
  "os"
  s "strings"
  // "sync"

  "github.com/bwmarrin/discordgo"
)

func errCheck(msg string, err error) {
  if err != nil {
    fmt.Printf("%s: %+v", msg, err)
  }
}

func CommandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
  user := message.Author
  botID := discord.State.User.ID
  if user.ID == botID || user.Bot {
    //Do nothing because the bot is talking
    return
  }

  commandPrefix := "&"
  content := message.Content


  if s.HasPrefix(content, commandPrefix + "hello") {
    discord.ChannelMessageSend(message.ChannelID, "Hello!")
  } else if s.HasPrefix(content, commandPrefix + "commands") {
    discord.ChannelMessageSend(message.ChannelID,
      "__**Command List**__\n" +
      "`hello: Test Bot returns a greeting`\n" +
      "`oof: Test Bot replies with a big oof`")
  } else if s.HasPrefix(content, commandPrefix + "oof") {
    f, err := os.Open("commands/pics/oof.png")
    if err != nil {
      errCheck("Something went wrong. Unable to open oof file at this time", err)
    } else {
      defer f.Close()
    }
    ms := &discordgo.MessageSend{
      Files: []*discordgo.File{
        &discordgo.File{
          Name:   "commands/pics/oof.png",
          Reader: f,
        },
      },
    }
    discord.ChannelMessageSendComplex(message.ChannelID, ms)
  } else if s.HasPrefix(content, commandPrefix + "rps") {
    rpsHandler(discord, message)
  }

  fmt.Printf("Message: %+v || From: %s\n\n", message.Message, message.Author)
}

func rpsHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
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

  dmChannels := make(chan *discordgo.Channel)
  // TODO: Add timeout for player choices

  go rpsContactPlayers(discord, message, dmChannels)

  for i := 0; i < 2; i++ {
    ch := <- dmChannels
    fmt.Println(ch.Name)
  }
}

func rpsContactPlayers(discord *discordgo.Session, message *discordgo.MessageCreate, dmChannels chan *discordgo.Channel) {
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

  dmChannels <- dmP1
  dmChannels <- dmP2

}
