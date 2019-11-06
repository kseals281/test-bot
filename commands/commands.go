package commands

import (
	"fmt"
	"os"
	s "strings"

	"github.com/bwmarrin/discordgo"
)

func errCheck(msg string, err error) {
	if err != nil {
		_ = fmt.Errorf("%s: %+v", msg, err)
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

	switch {

	case s.HasPrefix(content, commandPrefix+"hello"):
		_, err := discord.ChannelMessageSend(message.ChannelID, "Hello!")
		errCheck("", err)

	case s.HasPrefix(content, commandPrefix+"commands"):
		_, err := discord.ChannelMessageSend(message.ChannelID,
			"__**Command List**__\n"+
				"`hello: Test Bot returns a greeting`\n"+
				"`oof: Test Bot replies with a big oof`")
		errCheck("Failed to send list of commands", err)

	case s.HasPrefix(content, commandPrefix+"oof"):
		f, err := os.Open("commands/pics/oof.png")
		if err != nil {
			errCheck("Something went wrong. Unable to open oof file at this time", err)
		} else {
			defer f.Close()
		}
		ms := &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "commands/pics/oof.png",
					Reader: f,
				},
			},
		}
		_, err = discord.ChannelMessageSendComplex(message.ChannelID, ms)
		if err != nil {
			errCheck("Unable to send oof to channel", err)
		}

	case s.HasPrefix(content, commandPrefix+"rps"):
		RPSHandler(discord, message)
	}

	fmt.Printf("Message: %+v || From: %s\n\n", message.Message, message.Author)
}
