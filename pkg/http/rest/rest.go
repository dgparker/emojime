package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgparker/emojime/pkg/emojime"
	"github.com/fiseo/httpsrv"
	"github.com/husobee/vestigo"
	"go.uber.org/zap"
)

// EmojiService implementation provides access to an emoji database
type EmojiService interface {
	ListEmojis() ([]*emojime.Emoji, error)
	GetEmoji(name string) (*emojime.Emoji, error)
	SearchEmojis(query string) ([]*emojime.Emoji, error)
}

// New creates a new emojime rest server
func New(logger *zap.Logger, svc EmojiService) *http.Server {
	return httpsrv.NewWithDefault(newRouter(logger, svc))
}

func newRouter(logger *zap.Logger, svc EmojiService) http.Handler {
	r := vestigo.NewRouter()

	r.Get("/healthz", healthzHandler(logger))
	r.Get("/emojis", listEmojisHandler(logger, svc))
	r.Get("/emojis/:name", getEmojiHandler(logger, svc))

	return r
}

func healthzHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st := logRequest(logger, r, "healthz")

		logSummary(logger, r, "healthz", st)
		w.WriteHeader(http.StatusOK)
	}
}

func listEmojisHandler(logger *zap.Logger, svc EmojiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st := logRequest(logger, r, "list emojis")

		var res []*emojime.Emoji
		var err error

		query := r.URL.Query().Get("s")
		if query != "" {
			res, err = svc.SearchEmojis(query)
			if err != nil {
				logger.Error("list emojis handler search error", zap.String("error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res, err = svc.ListEmojis()
			if err != nil {
				logger.Error("list emojis handler list error", zap.String("error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if len(res) == 0 || res == nil {
			res = make([]*emojime.Emoji, 0)
		}

		data, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logSummary(logger, r, "list emojis", st)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func getEmojiHandler(logger *zap.Logger, svc EmojiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st := logRequest(logger, r, "get emoji")

		res, err := svc.GetEmoji(vestigo.Param(r, "name"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		data, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logSummary(logger, r, "get emoji", st)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func logRequest(logger *zap.Logger, r *http.Request, handler string) (st time.Time) {
	st = time.Now()
	logger.Info(
		"request",
		zap.String("path", r.URL.Path),
		zap.String("handler", handler),
		zap.String("request-id", r.Header.Get("request-id")),
		zap.String("time", st.UTC().String()),
	)
	return st
}

func logSummary(logger *zap.Logger, r *http.Request, handler string, st time.Time) {
	logger.Info(
		"request summary",
		zap.String("path", r.URL.Path),
		zap.String("handler", handler),
		zap.String("request-id", r.Header.Get("request-id")),
		zap.Float64("time", time.Since(st).Seconds()),
	)
}
