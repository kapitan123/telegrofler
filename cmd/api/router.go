package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/internal/command"
	"github.com/kapitan123/telegrofler/internal/command/choosePidor"
)

// AK TODO we just can't pass all stuff here, we still need an abstraction to group configuration
// temp solution with direct handler function
func setupRouter(r *mux.Router, runner *command.Runner, pdr *choosePidor.ChoosePidor) {
	r.HandleFunc("/callback", messageHandler(runner)).Methods("POST")
	r.HandleFunc("/chat/{chatid}/{offx}/{offy}/pidoroftheday", choosePidorHandler(pdr)).Methods("POST")
}

// AK TODO send messages to a dead message quee
func messageHandler(runner *command.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.WithError(err).Error("Failed to decode the callback message.")

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logContent(update.Message)

		err = runner.Run(r.Context(), update.Message)
		if err != nil {
			log.WithError(err).Error("Failed trying to invoke a command.")
		}
		// Intentionally swallows all exception so messages are not resend
		w.WriteHeader(http.StatusOK)
	}
}

func choosePidorHandler(pdr *choosePidor.ChoosePidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatarg := mux.Vars(r)["chatid"]

		chatId, err := strconv.ParseInt(chatarg, 10, 64)

		offxarg := mux.Vars(r)["offx"]
		offx, _ := strconv.ParseInt(offxarg, 10, 64)

		offyarg := mux.Vars(r)["offy"]
		offy, _ := strconv.ParseInt(offyarg, 10, 64)
		if err != nil {
			log.WithError(err).Error("Failed trying to invoke a command.", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		err = pdr.ChoosePidorWithOffset(r.Context(), chatId, int(offx), int(offy))

		if err != nil {
			log.Error("Failed trying to invoke a command.", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func logContent(update *tgbotapi.Message) {
	ujs, _ := json.Marshal(update)
	log.Info("Callback content:", string(ujs))
}
