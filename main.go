package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	DiscordToken string
	PubgApiToken string
	Players      PlayerList
)

const BOT_PREFIX = "!manco"
const PUBG_API_URL = "https://api.pubg.com/shards/steam"

func main() {
	log.Println("Manco Stats Bot Starting")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DiscordToken = os.Getenv("DISCORD_TOKEN")
	if DiscordToken == "" {
		log.Fatal("DISCORD_TOKEN not set")
	}
	PubgApiToken = os.Getenv("PUBG_API_TOKEN")
	if PubgApiToken == "" {
		log.Fatal("PUBG_API_TOKEN not set")
	}

	Players, err = Players.ReadPlayersCSV("players.csv")
	if err != nil {
		log.Fatal(err)
	}
	if len(Players) < 1 {
		log.Println("[WARN] No players loaded from CSV")
	}

	discordClient, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	discordClient.AddHandler(messageCreateHandler)
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages
	discordClient.Identify.Intents |= discordgo.IntentMessageContent

	err = discordClient.Open()
	if err != nil {
		log.Fatalf("Error opening Discord connection: %v", err)
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discordClient.Close()
}

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // ignore self messages
		return
	}

	if !strings.HasPrefix(m.Content, BOT_PREFIX) {
		return
	}

	arguments := strings.Split(m.Content, " ")
	if len(arguments) < 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Se debe enviar un comando")
		if err != nil {
			fmt.Println(err)
		}
	}

	command := arguments[1]

	switch command {
	case "help":
		SendHelpMessage(s, m)
	case "season":
		SendStats(Players, arguments[1:], s, m)
	case "playerlist":
		SendSavedPlayers(Players, s, m)
	case "saveplayer":
		SavePlayer(arguments[1:], s, m)
	default:
		SendUnrecognizedCommandMessage(s, m)
	}
}
