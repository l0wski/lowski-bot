package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const prefix string = "!lowski"

func main() {

	// grab bot token
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")

	// create session
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(messageCreate)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged // hmm

	err = sess.Open() // open session
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close() // close session after termination

	fmt.Println("bot is online brother.")

	// keep bot running until os interrupt (CTRL + C)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	case cmd == "play": // will likely have to break logic out to keep clean
		if len(args) >= 3 { // verify there is a url argument -- args[2]
			url := args[2]
			if isValidURL(url) {
				title, err := download(url)
				if err != nil {
					panic(err)
				}
				message := "Now Playing: " + title
				s.ChannelMessageSend(m.ChannelID, message)
			} else {
				s.ChannelMessageSend(m.ChannelID, "invalid url")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "you have to provide a link for me to play the music")
		}
	}
}

// download downloads the specified file
func download(url string) (filename string, err error) {
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "opus", url, "--print", "after_move:title")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	return string(out), nil
}

// isValidURL verifies url is valid
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
