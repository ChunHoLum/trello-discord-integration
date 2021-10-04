package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/ChunHoLum/trello-discord-bot/lib/logger"
)

type App struct {
	ObserverList []Observer
	conf         Config
}

func NewApp(conf Config) (*App, error) {
	app := &App{conf: conf}
	return app, nil
}

func (a *App) Run(ctx context.Context) error {

	go NewWebhookServer(ctx, a.conf.Server, a.onTrelloWebhook)

	const botCount = 1

	botList := [botCount]string{
		"discord",
		// "discord",
	}

	var doneChs [botCount]chan bool
	var botChs [botCount]chan os.Signal

	for i := range botList {
		botChs[i] = make(chan os.Signal, 1)
		doneChs[i] = make(chan bool)
		NewBot(ctx, a.SubscribeTrelloRequest, botList[i], a.conf, botChs[i], doneChs[i], uint64(i))
		signal.Notify(botChs[i], syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	}

	allTrue := false

Exit:
	for {
		if allTrue {
			break
		}
		resolveList := [botCount]bool{}
		for i := range botList {
			s := <-doneChs[i]
			resolveList[i] = s
		}
		for i := range resolveList {
			if !resolveList[i] {
				allTrue = false
				continue Exit
			}
			allTrue = true
		}
	}
	return nil

}

func (a *App) onTrelloWebhook(ctx context.Context, webhook Webhook) error {
	log := logger.Get(ctx)

	eventTranslationKey := webhook.Action.Display.TranslationKey

	if eventTranslationKey == "unknown" {
		eventTranslationKey = webhook.Action.Type
		webhook.Action.Display.TranslationKey = eventTranslationKey
	}

	log.Debugf("Event %s received", eventTranslationKey)
	for _, Observer := range a.ObserverList {
		Observer.update(webhook)
	}

	return nil
}

func (a *App) SubscribeTrelloRequest(o Observer) {
	a.ObserverList = append(a.ObserverList, o)
}

func (a *App) UnSubscribeTrelloRequest(o Observer) {
	a.ObserverList = removeFromslice(a.ObserverList, o)
}

func removeFromslice(ObserverList []Observer, ObserverToRemove Observer) []Observer {
	ObserverListLength := len(ObserverList)
	for i, Observer := range ObserverList {
		if ObserverToRemove.getName() == Observer.getName() {
			ObserverList[ObserverListLength-1], ObserverList[i] = ObserverList[i], ObserverList[ObserverListLength-1]
			return ObserverList[:ObserverListLength-1]
		}
	}
	return ObserverList
}
