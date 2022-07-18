package choosePidor

import (
	"context"
	_ "embed"

	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/storage"

	"github.com/kapitan123/telegrofler/internal/faceFinder"
	log "github.com/sirupsen/logrus"
)

// AK TODO remove duplicate embedding
//go:embed pidormark.png
var pidormarkPicture []byte

//go:embed tinfoil.jpg
var tinfoilPicture []byte

//go:embed test.jpg
var alesha []byte

const commandName = "choosePidor"

type ChoosePidor struct {
	messenger   messenger
	storage     pidorStorage
	watermarker watermarker
	systemclock systemclock
}

type (
	watermarker interface {
		Apply(bakground []byte, foreground []byte) ([]byte, error)
		ApplyOnAxis(bakground []byte, foreground []byte, offsetx int, offsety int, lenght int) ([]byte, error)
	}

	messenger interface {
		SendText(chatId int64, text string) error
		SendImg(chatId int64, img []byte, imgName string, caption string) error
		GetChatAdmins(chatId int64) ([]tgbotapi.ChatMember, error)
		GetUserCurrentProfilePic(userId int64) ([]byte, error)
	}

	pidorStorage interface {
		GetPidorForDate(ctx context.Context, chatid int64, date time.Time) (storage.Pidor, bool, error)
		CreatePidor(ctx context.Context, chatid int64, username string, userId int64, date time.Time) error
	}

	systemclock interface {
		Now() time.Time
	}
)

func New(messenger messenger, storage pidorStorage, watermarker watermarker, systemclock systemclock) *ChoosePidor {
	return &ChoosePidor{
		messenger:   messenger,
		storage:     storage,
		watermarker: watermarker,
		systemclock: systemclock,
	}
}

func (h *ChoosePidor) Handle(ctx context.Context, m *tgbotapi.Message) error {
	return h.ChoosePidor(ctx, m.Chat.ID)
}

func (h *ChoosePidor) ChoosePidor(ctx context.Context, chatId int64) error {
	currDate := h.systemclock.Now()
	pidor, found, err := h.storage.GetPidorForDate(ctx, chatId, currDate)

	if err != nil {
		return err
	}

	if found {
		err = h.messenger.SendText(chatId, pidor.UserName+" is still sucking juicy cocks")
		return err
	}

	admins, err := h.messenger.GetChatAdmins(chatId)

	if err != nil {
		return err
	}

	chosenOne := chooseRandom(admins)

	err = h.storage.CreatePidor(ctx, chatId, chosenOne.User.UserName, chosenOne.User.ID, currDate)
	if err != nil {
		return err
	}

	ppic, err := h.messenger.GetUserCurrentProfilePic(chosenOne.User.ID)

	if err != nil {
		log.WithError(err).Error("failed to generate user profile pic")
		return h.messenger.SendImg(chatId, tinfoilPicture, "tinfoil.png", "Скрытный пидор дня у нас @"+chosenOne.User.UserName)
	}

	markedPic, err := h.watermarker.Apply(ppic, pidormarkPicture)

	if err != nil {
		return err
	}

	return h.messenger.SendImg(chatId, markedPic, "pidor.png", "Pidor of the day is "+chosenOne.User.UserName)
}

func chooseRandom(memebers []tgbotapi.ChatMember) tgbotapi.ChatMember {
	return memebers[rand.Intn(len(memebers)-1)]
}

func (h *ChoosePidor) ShouldRun(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == commandName
}

func (h *ChoosePidor) ChoosePidorWithOffset(ctx context.Context, chatId int64, offx int, offy int) error {
	x, y, length, err := faceFinder.FindMouth()
	//x, y := 262.0814, 264.62463, 77.07892

	ppic := alesha
	markedPic, err := h.watermarker.ApplyOnAxis(ppic, pidormarkPicture, int(x)+offx, int(y)+offy, int(length))

	if err != nil {
		return err
	}

	return h.messenger.SendImg(chatId, markedPic, "pidor.png", "Pidor of the day is ") //+chosenOne.User.UserName)
}
