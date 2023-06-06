package port

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/common/server/httperr"
	"github.com/kapitan123/telegrofler/service/bot/app"
	"github.com/sirupsen/logrus"
)

// implements interface generated by openapi spec
type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app: app}
}

func (h HttpServer) HandleVideoSavedMessage(w http.ResponseWriter, r *http.Request) {
	postMessage := PubSubMessage{}
	if err := render.Decode(r, &postMessage); err != nil {
		httperr.BadRequest("invalid format of pubsub message", err, w, r)
		return
	}

	var savedVideoMessage VideoSavedMessage
	decodedDataBytes, err := base64.StdEncoding.DecodeString(postMessage.Message.Data)
	if err != nil {
		httperr.BadRequest("invalid encoding of pubsub data", err, w, r)
		return
	}

	err = json.Unmarshal(decodedDataBytes, &savedVideoMessage)
	if err != nil {
		httperr.BadRequest("invalid pubsub.Data payload", err, w, r)
		return
	}

	err = h.app.PublishVideo(r.Context(), savedVideoMessage.OriginalUrl, savedVideoMessage.SavedVideoAddr)

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Status(r, 200)
}

func (h HttpServer) HandleTelegramMessage(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		httperr.BadRequest("invalid format of telegram update message", err, w, r)
		return
	}

	LogBody(update)

	err = h.app.HandleTelegramMessage(r.Context(), update.Message)

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Status(r, 200)
}

func (h HttpServer) HandleVideoSaveFailedMessage(w http.ResponseWriter, r *http.Request) {
	// AK TODO not implemented

	render.Status(r, 200)
}

type VideoSavedMessage struct {
	SavedVideoAddr string `json:"saved_video_addr"`
	OriginalUrl    string `json:"original_url"`
}

// AK TODO TEMP to get raw result
func LogBody(upd tgbotapi.Update) {
	jsonBytes, _ := json.Marshal(upd)

	log := logrus.WithField("body", string(jsonBytes))

	log.Info("Tg callback recived")
}
