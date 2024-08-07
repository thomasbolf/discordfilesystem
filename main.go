package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {

	uploadFlag := flag.String("u", "", "Upload a file")
	downloadFlag := flag.String("d", "", "Download a file")
	flag.Parse()
	if *uploadFlag != "" {
		fmt.Println("Uploading file: ", *uploadFlag)
		upload(*uploadFlag)
		//testing download... this is a placeholder that will be removed
		//id 1: 1270509945401643008
		//id 2: 1270509946647347264
		token := os.Getenv("TOKEN")
		dg, _ := discordgo.New("Bot " + token)
		fmt.Println("about to download")
		res, err := dg.ChannelMessage("1268969289373843456", "1270509945401643008")
		if err != nil {
			panic("Error reading message")
		}
		fmt.Println("downloaded")
		resURL := res.Attachments[0].URL
		resp, err := http.Get(resURL)
		if err != nil {
			panic("Error getting HTTP response with attachment url")
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic("Error reading body")
		}
		resString := string(body)

		fmt.Println(resString)
		return
	}
	if *downloadFlag != "" {
		fmt.Println("Downloading file: ", *downloadFlag)
		download(*downloadFlag)
		return
	}
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	//error check to make sure the token is provided
	if token == "" {
		panic("No token provided")
	}
	//create a new discord bot session
	dg, err := discordgo.New("Bot " + token)
	//being really careful here...
	if err != nil {
		panic("Error creating discord session")
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		panic("Error opening connection")
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	fmt.Println(m.ChannelID)
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

}
