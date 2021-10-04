package main

import (
	"context"

	logger "github.com/ChunHoLum/trello-discord-bot/lib/logger"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Log     logger.Config `toml:"log"`
	Server  ServerConfig  `toml:"server"`
	Discord DiscordConfg  `toml:"discord"`
}

type ServerConfig struct {
	Cert string `toml:"cert"`
	Key  string `toml:"key"`
}

type DiscordConfg struct {
	Token     string `toml:"token"`
	ChannelId string `toml:"channelId"`
}

type SubscritionFunc func(o Observer)
type WebhookFunc func(ctx context.Context, webhook Webhook) error

type WebhookServer struct {
	onWebhook WebhookFunc
	counter   uint64
}

type Webhook struct {
	Model  Model  `json:"model"`
	Action Action `json:"action"`
}

type Model struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ShortUrl string `json:"shortUrl"`
	Perf     Perf   `json:"prefs"`
}

type Perf struct {
	BackgroundColor string `json:"backgroundColor"`
}

type Action struct {
	Type          string        `json:"type"`
	Data          ActionData    `json:"data"`
	Date          string        `json:"date"`
	Display       Display       `json:"display"`
	Member        Member        `json:"member"`
	MemberCreator MemberCreator `json:"memberCreator"`
}

type Display struct {
	TranslationKey string `json:"translationKey"`
}

type ActionData struct {
	IdMember   string          `json:"idMember,omitempty"`
	Card       ActionDataCard  `json:"card,omitempty"`
	List       ActionDataList  `json:"list,omitempty"`
	Board      ActionDataBoard `json:"board,omitempty"`
	Text       string          `json:"text,omitempty"`
	ListBefore ListBefore      `json:"listBefore"`
	ListAfter  ListAfter       `json:"listAfter"`
}

type ListBefore struct {
	Name string `json:"name"`
}

type ListAfter struct {
	Name string `json:"name"`
}

type ActionDataCard struct {
	Id        string `json:"id,omitempty"`
	Desc      string `json:"desc,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortLink string `json:"shortLink,omitempty"`
}

type ActionDataList struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ActionDataBoard struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortLink string `json:"shortLink,omitempty"`
}

type Member struct {
	Id        string `json:"id"`
	AvatarUrl string `json:"avatarUrl"`
	FullName  string `json:"fullName"`
	Username  string `json:"username"`
}

type MemberCreator struct {
	Id        string `json:"id"`
	AvatarUrl string `json:"avatarUrl"`
	FullName  string `json:"fullName"`
	Username  string `json:"username"`
}

// type MessageEmbed struct {
// 	URL string `json:""`
// 	Title string `json:""`
// 	Description string `json:""`
// 	 string `json:""`
// 	URL string `json:""`
// }

type Bot struct {
	name   string
	logger log.FieldLogger
}

type Discord struct {
	*Bot
	conf DiscordConfg
	dg   *discordgo.Session
}

type Observer interface {
	update(Webhook)
	getName() string
}
