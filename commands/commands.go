package commands

import (
  "fmt"
  "os"

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

  switch content {

    case commandPrefix + "hello":
      discord.ChannelMessageSend(message.ChannelID, "Hello!")

    case commandPrefix + "commands":
      discord.ChannelMessageSend(message.ChannelID,
        "__**Command List**__\n" +
        "`hello: Test Bot returns a greeting`\n" +
        "`gbfc: Test Bot tells you about the GBFC Overwatch team`")

    case commandPrefix + "oof":
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

    default: {} // Do nothing

    }

  fmt.Printf("Message: %+v || From: %s\n\n", message.Message, message.Author)
}

func