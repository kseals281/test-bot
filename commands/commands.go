package commands

import (
  "fmt"
  "os"
  s "strings"

  "github.com/bwmarrin/discordgo"
)

// type rpsResponse struct {
//   choice discordgo.Emoji
//   timeout bool
// }

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
  } else if s.HasPrefix(content, commandPrefix + "test") {
    rpsMessageOpponent(discord, message.Mentions[0], message.Author.Username)
  }

  fmt.Printf("Message: %+v || From: %s\n\n", message.Message, message.Author)
}

func rpsHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
  // Webhook messages do not contain a full author (i.e. non-user author) and
  // should be ignored as a result
  if message.WebhookID != "" {
    return
  }

  // if message.MentionEveryone {
  //   discord.ChannelMessageSend(message.ChannelID,
  //     "You may not challenge everyone to a game of rock, paper, scissors")
  //   return
  // } else if len(message.MentionRoles) > 0 {
  //   discord.ChannelMessageSend(message.ChannelID,
  //     "You may not challenge a role to a game of rock, paper, scissors")
  //   return
  // } else if len(message.Mentions) > 1 {
  //   discord.ChannelMessageSend(message.ChannelID,
  //     "You can only challenge one person at a time to rock, paper, scissors")
  //   return
  // } else if len(message.Mentions) == 0 {
  //   discord.ChannelMessageSend(message.ChannelID,
  //     "You must select at least one opponent for rock, paper, scissors")
  //   return
  // }


  // opponent := message.Mentions[0]



  // go rpsMessageOpponent(discord, opponent, message.Author)
}

func rpsMessageOpponent(discord *discordgo.Session, opponent *discordgo.User, challenger string) {
  direct_message, err := discord.UserChannelCreate(opponent.ID)
  if err != nil {
    errCheck("Unable to create direct message with the opponent", err)
    return
  }

  rules := ("__**You have been challenged to a game of rock paper scissors!!!**__\n" +
           ":fist::raised_hand::v:\t:fist::raised_hand::v:\t:fist::raised_hand::v:\t" +
           ":fist::raised_hand::v:\t:fist::raised_hand::v:" +
           "\n\n" +
           challenger + " has challenged you to a match of rock paper scissors.")


  discord.ChannelMessageSend(direct_message.ID, rules)

}

// func testDirect(discord *discordgo.Session, message *discordgo.MessageCreate) {
//   opponent := message.Mentions[0]
//   direct_message, err := discord.UserChannelCreate(opponent.ID)

//   if err != nil {
//     fmt.Println("Error %v", err)
//     return
//   }

//   discord.ChannelMessageSend(direct_message.ID, "Test")

// }
