package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"./config"
	_ "github.com/lib/pq"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var db *sql.DB

//Quote ...
type Quote struct {
	ID     int
	Author string
	Body   string
}

func main() {
	var err error

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ", config.PgHost, config.PgPort, config.PgUser, config.PgPass, config.PgDB))
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	// fmt.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me Go or go for random quote")
			bot.Send(msg)
		case "Go":
			quote := new(Quote)
			quote = randomQuoteRun()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n(c)%s", quote.Body, quote.Author))
			bot.Send(msg)
		case "go":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Let`s go")
			bot.Send(msg)
		}
	}
}

func randomQuoteRun() *Quote {
	rows, err := db.Query("SELECT * FROM quotes")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	quotes := make([]*Quote, 0)
	for rows.Next() {
		quote := new(Quote)
		err := rows.Scan(&quote.ID, &quote.Author, &quote.Body)
		if err != nil {
			log.Fatal(err)
		}
		quotes = append(quotes, quote)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	myrand := random(0, len(quotes))
	randomQuote := quotes[myrand]
	return randomQuote
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	if min > max {
		return min
	}
	return rand.Intn(max-min) + min

}
