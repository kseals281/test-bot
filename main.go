package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"

  "test-bot/commands"

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
  errCheck("error retrieving account", err)

  botID = Session.State.User.ID
  Session.AddHandler(commands.CommandHandler)
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
  errCheck("Error opening connection to Discord", err)

  <-interrupt
}

func errCheck(msg string, err error) {
  if err != nil {
    fmt.Printf("%s: %+v", msg, err)
    panic(err)
  }
}

