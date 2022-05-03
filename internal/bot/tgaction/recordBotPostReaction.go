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
	reaction, err := bot.TryExtractVideoRepostReaction(m)
	if err != nil {
		return !isHandeled, err
	}

	if reaction.Sender == "" {
		return !isHandeled, nil
	}

	log.Infof("Reaction was found for %s sent by %s", reaction.VideoId, reaction.Sender)

	h.AddReactionToPost(reaction)

	return isHandeled, nil
}