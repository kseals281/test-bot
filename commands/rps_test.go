package commands

import (
	"github.com/bwmarrin/discordgo"
	"reflect"
	"testing"
)

func TestRPSHandler(t *testing.T) {
	type args struct {
		discord *discordgo.Session
		message *discordgo.MessageCreate
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_messageReaction(t *testing.T) {
	type args struct {
		reaction *discordgo.MessageReactions
		message  *discordgo.Message
	}
	tests := []struct {
		name string
		args args
		want *discordgo.MessageReaction
	}{
		{
			"One Reaction",
			args{
				&discordgo.MessageReactions{
					Count: 1,
					Me:    false,
					Emoji: nil,
				},
				&discordgo.Message{
					ID:        "",
					ChannelID: "",
					GuildID:   "",
					Author:    nil,
					Reactions: nil,
				},
			},
			&discordgo.MessageReaction{
				UserID:    "",
				MessageID: "",
				Emoji:     discordgo.Emoji{},
				ChannelID: "",
				GuildID:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := messageReaction(tt.args.reaction, tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("messageReaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rpsContactPlayers(t *testing.T) {
	type args struct {
		discord *discordgo.Session
		message *discordgo.MessageCreate
	}
	tests := []struct {
		name  string
		args  args
		want  *discordgo.Message
		want1 *discordgo.Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := rpsContactPlayers(tt.args.discord, tt.args.message)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rpsContactPlayers() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("rpsContactPlayers() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_rpsResults(t *testing.T) {
	tests := []struct {
		name       string
		rpsChoices [2]*roshambo
		want       string
	}{
		{
			"Rock & Rock",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✊"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✊"}}, &discordgo.User{ID: "Player 2"}},
			},
			"",
		}, {
			"Rock & Paper",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✊"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✋"}}, &discordgo.User{ID: "Player 2"}},
			},
			"Player 2",
		}, {
			"Rock & Scissors",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✊"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✌️"}}, &discordgo.User{ID: "Player 2"}},
			},
			"Player 1",
		}, {
			"Paper & Rock",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✋"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✊"}}, &discordgo.User{ID: "Player 2"}},
			},
			"Player 1",
		}, {
			"Paper & Paper",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✋"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✋"}}, &discordgo.User{ID: "Player 2"}},
			},
			"",
		}, {
			"Paper & Scissors",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✋"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✌️"}}, &discordgo.User{ID: "Player 2"}},
			},
			"Player 2",
		}, {
			"Scissors & Rock",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✌️"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✊"}}, &discordgo.User{ID: "Player 2"}},
			},
			"Player 2",
		}, {
			"Scissors & Paper",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✌️"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✋"}}, &discordgo.User{ID: "Player 2"}},
			},
			"Player 1",
		}, {
			"Scissors & Scissors",
			[2]*roshambo{
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✌️"}}, &discordgo.User{ID: "Player 1"}},
				{&discordgo.MessageReaction{Emoji: discordgo.Emoji{Name: "✌️"}}, &discordgo.User{ID: "Player 2"}},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rpsResults(tt.rpsChoices); got != tt.want {
				t.Errorf("rpsResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rpsWaitForReaction(t *testing.T) {
	type args struct {
		s         *discordgo.Session
		message   *discordgo.Message
		reactions chan *roshambo
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_rpsWinner(t *testing.T) {
	type args struct {
		discord    *discordgo.Session
		message    *discordgo.MessageCreate
		rpsChoices [2]*roshambo
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_aiChoice(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want *discordgo.MessageReaction
	}{
		{
			"Rock",
			0,
			&discordgo.MessageReaction{
				Emoji: discordgo.Emoji{Name: "✊"},
			},
		}, {
			"Paper",
			1,
			&discordgo.MessageReaction{
				Emoji: discordgo.Emoji{Name: "✋"},
			},
		}, {
			"Scissors",
			2,
			&discordgo.MessageReaction{
				Emoji: discordgo.Emoji{Name: "✌️"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aiChoice(tt.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("aiChoice() = %v, want %v", got, tt.want)
			}
		})
	}
}
