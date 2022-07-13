package replaceLinkWithMessage

import (
	"context"
	"time"

	"github.com/kapitan123/telegrofler/internal/contentLoader"
	"github.com/kapitan123/telegrofler/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type ReplaceLinkWithMessage struct {
	messenger  messenger
	storage    postStorage
	downloader downloader
}

type messenger interface {
	ReplyWithText(chatId int64, messageId int, text string) error
	SendTrackableVideo(chatId int64, linktToUserName string, trackToken string, title string, payload []byte) error
	Delete(chatId int64, messageId int) error
}

type postStorage interface {
	UpsertPost(ctx context.Context, p storage.Post) error
}

type downloader interface {
	DownloadContent(dUrl string) ([]byte, error)
	ExtractVideoMeta(url string) (*contentLoader.VideoMeta, error)
	CanExtractVideoMeta(url string) bool
}

func New(messenger messenger, storage postStorage, downloader downloader) *ReplaceLinkWithMessage {
	return &ReplaceLinkWithMessage{
		messenger:  messenger,
		storage:    storage,
		downloader: downloader,
	}
}

func (h *ReplaceLinkWithMessage) Handle(ctx context.Context, m *tgbotapi.Message) error {
	url, chatId, sender := m.Text, m.Chat.ID, m.From.UserName

	meta, err := h.downloader.ExtractVideoMeta(url)

	if err != nil {
		return err
	}

	log.Info("Url was found in a callback message: ", url)

	content, err := h.downloader.DownloadContent(url)

	if err != nil {
		return err
	}

	err = h.messenger.SendTrackableVideo(chatId, sender, meta.Id, meta.Title, content)

	if err != nil {
		return err
	}

	// we don't really care if if has failed and it makes integration tests a lot easier
	_ = h.messenger.Delete(chatId, m.MessageID)

	newPost := storage.Post{
		VideoId:        meta.Id,
		Source:         meta.Type,
		RoflerUserName: sender,
		Url:            url,
		Reactions:      []storage.Reaction{},
		KeyWords:       []string{},
		PostedOn:       time.Now(),
	}

	err = h.storage.UpsertPost(ctx, newPost)

	return err
}

func (h *ReplaceLinkWithMessage) ShouldRun(m *tgbotapi.Message) bool {
	return h.downloader.CanExtractVideoMeta(m.Text)
}
