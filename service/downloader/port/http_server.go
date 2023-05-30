// it is a wrapper around application to format requests
package port

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/kapitan123/telegrofler/common/server/httperr"
	"github.com/kapitan123/telegrofler/service/downloader/app"
)

// implements interface generated by openapi spec
type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app: app}
}

func (h HttpServer) HandleVideoUrlPublishedMessage(w http.ResponseWriter, r *http.Request) {
	postMessage := PostPubSubMessage{}
	if err := render.Decode(r, &postMessage); err != nil {
		httperr.BadRequest("invalid format of pubsub message", err, w, r)
		return
	}

	var videoMessage VideoUrlPublishedMessage
	err := json.Unmarshal(postMessage.Data, &videoMessage)
	if err != nil {
		httperr.BadRequest("invalid pubsub.Data payload", err, w, r)
		return
	}

	err = h.app.SaveVideoToStorage(r.Context(), videoMessage.Url)

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Status(r, 200)
}

type VideoUrlPublishedMessage struct {
	Url string `json:"url"`
}