package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {

	const prefix string = "!lowski"

	// grab bot token
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")

	// create session
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// if message author is bot itself, do nothing
		if m.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(m.Content, " ")

		if args[0] != prefix {
			return
		}

		// grab command
		cmd := args[1]
		switch {
		case cmd == "help":
			help := []string{
				"you seem like you need help a lot",
				"ever try helping yourself?",
				"let me hold a dollar",
				"lmao look at this guy asking for help",
				"not feeling helpful fr",
				"nah",
				"can't do it",
			}
			choice := rand.Intn(len(help))
			s.ChannelMessageSend(m.ChannelID, help[choice])
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged // hmm

	err = sess.Open() // open session
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close() // close session after termination

	fmt.Println("bot is online brother")


	// keep bot running until os interrupt (CTRL + C)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
