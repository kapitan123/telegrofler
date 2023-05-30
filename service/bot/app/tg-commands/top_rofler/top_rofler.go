package toprofler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/kapitan123/telegrofler/service/bot/domain/message"
	"github.com/kapitan123/telegrofler/service/bot/internal/messenger/format"
)

const commandName = "toprofler"

type TopRofler struct {
	messenger messenger
	storage   postStorage
}

type messenger interface {
	SendText(chatId int64, text string) (int, error)
}

type postStorage interface {
	GetChatPosts(ctx context.Context, chatId int64) ([]domain.Post, error)
}

func New(messenger messenger, storage postStorage) *TopRofler {
	return &TopRofler{
		messenger: messenger,
		storage:   storage,
	}
}

func (h *TopRofler) Handle(ctx context.Context, message message.Message) error {
	posts, err := h.storage.GetChatPosts(ctx, message.ChatId())
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return nil
	}

	// AK TODO also should be gone to domain
	roflerScores := countScores(posts)

	listMeassge := format.AsDescendingList(roflerScores, "🤡 <b>%s</b>: %d")

	_, err = h.messenger.SendText(message.ChatId(), listMeassge)
	if err != nil {
		return err
	}
	return nil
}

func countScores(posts []domain.Post) map[string]int {
	roflerScores := map[domain.UserRef]int{}
	for _, p := range posts {
		roflerScores[p.Poster] += len(p.Reactions)
	}

	names := map[string]int{}
	for k, v := range roflerScores {
		names[k.DisplayName] = v
	}
	return names
}

func (h *TopRofler) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}
