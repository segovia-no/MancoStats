package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"pubgstats/pubgDAL"
	"strings"
	"syscall"
)

var (
	DiscordToken string
	Servers      []ServerPlayerList
	BotPrefix    string
	PubgDAL      *pubgDAL.PUBGApiDAL
	pubgApiURL   string
	pubgApiToken string
)

const DefaultPUBGApiURL = "https://api.pubg.com/shards/steam"
const DefaultBotPrefix = "!manco"

func main() {
	log.Println("Manco Stats Bot Starting")

	err := loadEnvVars()
	if err != nil {
		log.Fatal(err)
	}

	PubgDAL, err = pubgDAL.NewPUBGApiDAL(pubgApiURL, pubgApiToken)
	if err != nil {
		log.Fatal(err)
	}

	err = LoadServerCSVs(&Servers)
	if err != nil {
		fmt.Println("[WARN] Couldn't load servers: " + err.Error())
	}

	discordClient, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	discordClient.Identify.Intents = discordgo.IntentsGuildMessages
	discordClient.Identify.Intents |= discordgo.IntentMessageContent

	discordClient.AddHandler(messageCreateHandler)

	err = discordClient.Open()
	if err != nil {
		log.Fatalf("Error opening Discord connection: %v", err)
	}

	log.Println("Bot is now running, bot prefix is: ", BotPrefix, " - Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discordClient.Close()
}

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // ignore self messages
		return
	}

	if !strings.HasPrefix(m.Content, BotPrefix) {
		return
	}

	arguments := strings.Split(m.Content, " ")
	if len(arguments) < 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Se debe enviar un comando")
		if err != nil {
			fmt.Println(err)
		}
	}

	srvIdx, err := GetServerIndex(m.GuildID, &Servers)
	if err != nil {
		fmt.Println(err)
		_, err = s.ChannelMessageSend(m.ChannelID, "No se pudo encontrar el servidor")
		return
	}

	command := arguments[1]
	switch command {
	case "help":
		SendHelpMessage(s, m)
	case "stats", "season", "weekly", "semana":
		SendSeasonStats(Servers[srvIdx].PlayerList, arguments[1:], s, m)
	case "addiction", "vicio":
		SendAddictionStats(Servers[srvIdx].PlayerList, s, m)
	case "playerlist":
		SendSavedPlayers(Servers[srvIdx].PlayerList, s, m)
	case "saveplayer":
		SavePlayer(arguments[1:], &Servers[srvIdx], s, m)
	case "removeplayer":
		RemovePlayer(arguments[1:], &Servers[srvIdx], s, m)
	default:
		SendUnrecognizedCommandMessage(s, m)
	}
}

func loadEnvVars() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	DiscordToken = os.Getenv("DISCORD_TOKEN")
	if DiscordToken == "" {
		return errors.New("DISCORD_TOKEN not set")
	}
	pubgApiToken = os.Getenv("PUBG_API_TOKEN")
	if pubgApiToken == "" {
		return errors.New("PUBG_API_TOKEN not set")
	}

	readBotPrefix := os.Getenv("BOT_PREFIX")
	if readBotPrefix == "" {
		BotPrefix = DefaultBotPrefix
	} else {
		BotPrefix = readBotPrefix
	}

	readPUBGApiURL := os.Getenv("PUBG_API_URL")
	if readPUBGApiURL == "" {
		pubgApiURL = DefaultPUBGApiURL
	} else {
		pubgApiURL = readPUBGApiURL
	}

	return nil
}
