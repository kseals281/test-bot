package commands

import (
	"fmt"

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

	commandPrefix := "~"
	content := message.Content

	switch content {

	case commandPrefix + "hello":
		_, err := discord.ChannelMessageSend(message.ChannelID, "Hello!")
		errCheck("Failed to send \"Hello!\" to channel", err)

	case commandPrefix + "commands":
		_, err := discord.ChannelMessageSend(message.ChannelID,
			"__**Command List**__\n"+
				"`hello: Test Bot returns a greeting`\n")
		errCheck("Failed to provide Command List to channel", err)

	default:
		{
		} // Do nothing

	}

	fmt.Printf("Message: %+v || From: %s\n\n", message.Message, message.Author)
}
