package main

import (
  "fmt"

  "github.com/bwmarrin/discordgo"
)

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
  user := message.Author
  if user.ID == botID || user.Bot {
    //Do nothing because the bot is talking
    return
  }

  content := message.Content

  switch content {
    case "!hello":
      discord.ChannelMessageSend(message.ChannelID, "Hello!")
    case "!gbfc":
      discord.ChannelMessageSend(
        message.ChannelID,
        "Genji Butt Fan Club aka a bunch of try hard casuals")
    case "!commands":
      discord.ChannelMessageSend(message.ChannelID,
        "__**Command List**__\n" +
        "`hello: GBSB returns a greeting`\n" +
        "`gbfc: GBSB tells you about the GBFC Overwatch team`")
    case "!oof":
      f, err := os.Open("oof.png")
      if err != nil {
        return nil, err
      }
      defer f.Close()

      ms := &discordgo.MessageSend{
        Embed: &discordgo.MessageEmbed{
          Image: &discordgo.MessageEmbedImage{
            URL: "attachment://" + fileName,
          },
        },
        Files: []*discordgo.File{
          &discordgo.File{
            Name:   fileName,
            Reader: f,
          },
        },
      }

      s.ChannelMessageSendComplex(channelID, ms)
    default: { // Only reply if the message was a command
      if (string(content[0]) == "!") {
        discord.ChannelMessageSend(message.ChannelID, "Command not found. List of commands comming soon.")
      }
    }
  }

  fmt.Printf("Message: %+v || From: %s\n\n", message.Message, message.Author)
}