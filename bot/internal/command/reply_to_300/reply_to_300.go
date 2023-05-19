package replyTo300

import (
	"context"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var regexp300 = regexp.MustCompile(`300|Триста|триста`)

type ReplyTo300 struct {
	messenger messenger
	queue     queue
}

type messenger interface {
	ReplyWithText(chatId int64, messageId int, text string) (int, error)
}

// AK TODO i might change signature of Handle method, and wrap callbacks in
// selfDestructable decorator
type queue interface {
	EnqueueDeleteMessage(chatId int64, messageId int) error
}

func New(messenger messenger, queue queue) *ReplyTo300 {
	return &ReplyTo300{
		messenger: messenger,
		queue:     queue,
	}
}

func (h *ReplyTo300) Handle(ctx context.Context, m *tgbotapi.Message) error {
	chatId := m.Chat.ID
	newMessageId, err := h.messenger.ReplyWithText(chatId, m.MessageID, "🤣🚜 ♂ Отсоси у тракториста ♂ 🚜🤣")

	if err != nil {
		return err
	}

	return h.queue.EnqueueDeleteMessage(chatId, newMessageId)
}

func (h *ReplyTo300) ShouldRun(m *tgbotapi.Message) bool {
	return regexp300.MatchString(m.Text)
}