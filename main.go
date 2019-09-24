package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"

  "github.com/bwmarrin/discordgo"
)

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session, _ = discordgo.New()

var (
    commandPrefix string
    botID         string
)

// Read in all options from environment variables and command line arguments.
func init() {

  // Discord Authentication Token
  Session.Token = os.Getenv("DISCORD_TEST_BOT")
  if Session.Token == "" {
    // Pointer, flag, default, description
    flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
  }
}

func main() {

  // Declare any variables needed later.
  var err error

  // Setup interrupt
  interrupt := make(chan os.Signal)
  signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

  // Parse command line arguments
  flag.Parse()

  // Verify a Token was provided
  if Session.Token == "" {
    log.Println("You must provide a Discord authentication token.")
    return
  }

  // Verify the Token is valid and grab user information
  Session.State.User, err = Session.User("@me")
  errCheck("error retrieving account", err, true)

  botID = Session.State.User.ID
  Session.AddHandler(commandHandler)
  Session.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
    err = discord.UpdateStatus(0, "A friendly development bot!")
    if err != nil {
      fmt.Println("Error attempting to set my status")
    }
    servers := discord.State.Guilds
    fmt.Printf("Test-Bot has started on %d servers", len(servers))
  })

  // Open a websocket connection to Discord
  err = Session.Open()
  defer Session.Close()
  errCheck("Error opening connection to Discord", err, true)

  commandPrefix = "&"

  <-interrupt
}

func errCheck(msg string, err error, die bool) {
  if err != nil {
    fmt.Printf("%s: %+v", msg, err)
    if die {
      panic(err)
    }
  }
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
  user := message.Author
  if user.ID == botID || user.Bot {
    //Do nothing because the bot is talking
    return
  }

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
        errCheck("Something went wrong. Unable to open oof file at this time", err, false)
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
