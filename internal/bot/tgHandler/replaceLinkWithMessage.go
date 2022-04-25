package tgHandler

import (
	"time"

	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data"
	"github.com/kapitan123/telegrofler/internal/data/model"
	"github.com/kapitan123/telegrofler/pkg/source/sourceFactory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type ReplaceLinkWithMessage struct {
	*bot.Bot
	*data.PostsStore
}

func (h *ReplaceLinkWithMessage) Handle(mess *tgbotapi.Message) (bool, error) {
	isHandeled := true
	source, found := sourceFactory.TryGetSource(mess.Text)
	if !found {
		return !isHandeled, nil
	}

	svp := h.ConvertToSourceVideoPost(mess)

	if svp == nil {
		return !isHandeled, nil
	}

	log.Info("Url was found in a callback message: ", svp.Url)

	evi, err := source.ExtractVideoFromUrl(svp.Url)

	if err != nil {
		return isHandeled, err
	}

	svp.VideoData.Payload = evi.Payload
	svp.VideoData.Title = evi.Title
	svp.VideoData.Id = evi.Id

	log.Info("Trying to post to telegram: ", svp)

	err = h.RepostConvertedVideo(svp)

	if err != nil {
		return isHandeled, err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = h.DeletePost(svp.ChatId, svp.OriginalMessageId)

	newPost := model.Post{
		VideoId:        svp.VideoData.Id,
		Source:         evi.Type,
		RoflerUserName: svp.Sender,
		Url:            svp.Url,
		Reactions:      []model.Reaction{},
		KeyWords:       []string{},
		PostedOn:       time.Now(),
	}

	h.Upsert(newPost)

	if err != nil {
		return isHandeled, err
	}

	return isHandeled, nil
}
