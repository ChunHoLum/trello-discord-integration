package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

func NewBot(
	ctx context.Context,
	subscribe SubscritionFunc,
	botName string,
	conf Config,
	botChan chan os.Signal,
	doneCh chan bool,
	counter uint64,
) {

	switch botName {
	case "discord":
		botId := fmt.Sprintf("%v-%v", time.Now().Unix(), atomic.AddUint64(&counter, 1))

		bot := make(chan *Discord)
		go NewDiscordBot(ctx, conf.Discord, botChan, doneCh, bot, botId)
		b := <-bot
		subscribe(b)

	default:
		fmt.Println("BOT NOT FOUND")
	}
}

func (bot Bot) InitBot() string {

	// returns comic universe
	return bot.name
}
