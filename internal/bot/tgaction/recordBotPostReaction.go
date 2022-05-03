package tgaction

import (
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/firestore"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type RecordBotPostReaction struct {
	*bot.Bot
	*firestore.PostsStore
}

func NewRecordBotPostReaction(b *bot.Bot, ps *firestore.PostsStore) *RecordBotPostReaction {
	return &RecordBotPostReaction{
		Bot:        b,
		PostsStore: ps,
	}
}

func (h *RecordBotPostReaction) Handle(m *tgbotapi.Message) (bool, error) {
	isHandeled := true
	mediaRepy, err := bot.TryExtractVideoRepostReaction(m)
	reaction := mediaRepy.Reaction

	if err != nil {
		return !isHandeled, err
	}

	if reaction.Sender == "" {
		return !isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", mediaRepy.VideoId, reaction.Sender)

	exPost, found, err := h.GetById(mediaRepy.VideoId)

	// in this case we don't record reaction as all bot posts should be saved already
	if !found {
		return isHandeled, nil
	}

	if err != nil {
		return isHandeled, err
	}

	exPost.AddReaction(reaction.Sender, reaction.Text, reaction.MessageId)
	h.Upsert(exPost)
	return isHandeled, nil
}
