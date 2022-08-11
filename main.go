package main

import (
	"database/sql"
	"flag"

	"github.com/AyJayTee/emotitron-go/bot"
	"github.com/AyJayTee/emotitron-go/components"
)

var (
	Token string
	db    *sql.DB
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	db = components.StartDatabase()
	defer db.Close()

	bot.Start(Token, db)
}
