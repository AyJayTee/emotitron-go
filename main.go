package main

import (
	"flag"

	"github.com/AyJayTee/emotitron-go/bot"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	bot.Start(Token)
}
