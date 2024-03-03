package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/LigeronAhill/friday_bot/clients/telegram"
	event_consumer "github.com/LigeronAhill/friday_bot/consumer/event-consumer"
	ev "github.com/LigeronAhill/friday_bot/events/telegram"

	// "github.com/LigeronAhill/friday_bot/storage/files"
	"github.com/LigeronAhill/friday_bot/storage/pg"
	_ "github.com/lib/pq"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	token := os.Getenv("TELEGRAM_APITOKEN")
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbname := os.Getenv("PGDATABASE")
	sslmode := os.Getenv("PGSSLMODE")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	s, err := pg.New(psqlInfo)
	if err != nil {
		log.Fatal("can't connect to db", err)
	}
	if err := s.Init(); err != nil {
		log.Fatal("can't init storage: ", err)
	}
	tgClient := telegram.New(tgBotHost, token)
	// tgClient := telegram.New(tgBotHost, mustToken())

	eventsProcessor := ev.New(tgClient, s)

	log.Print("service started")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	// bot -tg-bot-token 'my token'
	token := flag.String("t", "", "token for access telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
