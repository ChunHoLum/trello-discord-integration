package main

import (
	"context"
	"errors"
	"os"
	"strconv"
	"unicode/utf8"

	logger "github.com/ChunHoLum/trello-discord-bot/lib/logger"
	"github.com/bwmarrin/discordgo"
)

func NewDiscordBot(
	ctx context.Context,
	conf DiscordConfg,
	discordChan chan os.Signal,
	doneCh chan bool,
	bot chan *Discord,
	botId string,
) {

	ctx, _ = logger.WithField(ctx, "bot", "discord"+botId)
	log := logger.Get(ctx)

	log.Info("Starting discord bot...")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Errorf("error creating Discord session,", err)
		// return nil, err
	}
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Errorf("error opening connection,", err)
		// return nil, err
	}

	discord := &Discord{
		conf: conf,
		dg:   dg,
		Bot: &Bot{
			name:   "discord",
			logger: log,
		},
	}

	bot <- discord
	log.Info("Discord Bot Started")

	// dg.ChannelMessageSendEmbed()
	<-discordChan
	log.Errorf("Gracefully  discord session")
	discord.dg.Close()
	doneCh <- true

}
func (d *Discord) getName() string {
	return d.name
}

func (d *Discord) update(webhook Webhook) {
	d.logger.Infof("Sending message to discord channel %s.", d.conf.ChannelId)
	emb, err := d.generateEmbed(webhook)
	if err != nil {
		d.logger.Warn(err)
		return
	}
	d.dg.ChannelMessageSendEmbed(d.conf.ChannelId, &emb)

	// d.dg.ChannelMessageSend(d.conf.ChannelId, webhook.Action.Display.TranslationKey)
}

func (d *Discord) generateEmbed(webhook Webhook) (discordgo.MessageEmbed, error) {

	color := webhook.Model.Perf.BackgroundColor
	colorInt, err := strconv.ParseInt(TrimFirstRune(color), 16, 32)
	if err != nil {
		d.logger.Error(err)
		return discordgo.MessageEmbed{}, err
	}

	author := &discordgo.MessageEmbedAuthor{
		Name:    webhook.Action.MemberCreator.FullName,
		IconURL: webhook.Action.MemberCreator.AvatarUrl + "/50.png",
	}

	title, err := d.generateEmbedTitle(webhook)
	if err != nil {
		return discordgo.MessageEmbed{}, err
	}
	return discordgo.MessageEmbed{
		URL:         "https://trello.com/c/rZvmyavS" + webhook.Action.Data.Card.ShortLink,
		Title:       title,
		Timestamp:   webhook.Action.Date,
		Description: webhook.Action.Data.Card.Desc,
		Color:       int(colorInt),
		Author:      author,
	}, nil
}

const (
	EVENT_ADD_ATTACHMENT_TO_CARD = "action_add_attachment_to_card"
	EVENT_ADD_CHECKLIST_TO_CARD  = "action_add_checklist_to_card"
	// EVENT_ADD_LABEL_TO_CARD           = "action_add_label_to_card"
	EVENT_ADD_MEMBER_TO_CARD          = "action_member_joined_card"
	EVENT_COMMENT_CARD                = "action_comment_on_card"
	EVENT_DELETE_ATTACHMENT_FROM_CARD = "action_delete_attachment_from_card"
	EVENT_DELETE_CARD                 = "action_delete_card"
	EVENT_ARCHIVE_CARD                = "action_archived_card"
	EVENT_DELETE_COMMENT              = "deleteComment" // exclude type in trello, returning unknown
	EVENT_CREATE_CARD                 = "action_create_card"
	EVENT_REMOVE_CHECKLIST_FROM_CARD  = "action_remove_checklist_from_card"
	// EVENT_REMOVE_LABEL_FROM_CARD      = "action_remove_label_from_card"
	EVENT_REMOVE_MEMBER_FROM_CARD = "action_member_left_card"
	EVENT_UPDATE_CHECKITEM_STATE  = "action_completed_checkitem"
	EVENT_MOVE_CARD_LIST_TO_LIST  = "action_move_card_from_list_to_list"
)

func (d *Discord) generateEmbedTitle(webhook Webhook) (string, error) {

	switch webhook.Action.Display.TranslationKey {

	case EVENT_ADD_ATTACHMENT_TO_CARD:
		return "Attachment added to card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_ADD_CHECKLIST_TO_CARD:
		return "CheckList added to card:" + webhook.Action.Data.Card.Name, nil
	case EVENT_ADD_MEMBER_TO_CARD:
		return "Member added to card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_COMMENT_CARD:
		return "New comment on card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_DELETE_ATTACHMENT_FROM_CARD:
		return "Attachment deleted from card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_DELETE_CARD:
		return "Delete card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_ARCHIVE_CARD:
		return "Archived card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_DELETE_COMMENT:
		return "Comment delete on card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_CREATE_CARD:
		return "Created card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_REMOVE_CHECKLIST_FROM_CARD:
		return "CheckList removed from card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_REMOVE_MEMBER_FROM_CARD:
		return "Member left card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_UPDATE_CHECKITEM_STATE:
		return "CheckItem updated on card: " + webhook.Action.Data.Card.Name, nil
	case EVENT_MOVE_CARD_LIST_TO_LIST:
		return "Card: " + webhook.Action.Data.Card.Name + " moved from " + webhook.Action.Data.ListBefore.Name + " to " + webhook.Action.Data.ListAfter.Name, nil
	default:
		return "", errors.New("Event " + webhook.Action.Display.TranslationKey + " currently not supported")
	}

}
func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
