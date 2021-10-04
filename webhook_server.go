package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ChunHoLum/trello-discord-bot/lib"
	logger "github.com/ChunHoLum/trello-discord-bot/lib/logger"
)

const (
	addr = "0.0.0.0"
	port = "443"
)

func NewWebhookServer(ctx context.Context, conf ServerConfig, onWebhook WebhookFunc) error {

	ctx, _ = logger.WithField(ctx, "addr", addr)
	ctx, _ = logger.WithField(ctx, "port", port)

	log := logger.Get(ctx)

	log.Info("Start trello webhook server")
	srv := &WebhookServer{
		onWebhook: onWebhook,
	}

	http.HandleFunc("/", srv.processWebhook)
	err := http.ListenAndServeTLS(strings.Join([]string{addr, port}, ":"), conf.Cert, conf.Key, nil)

	if err != nil {
		logger.Standard().Fatal("ListenAndServe: ", err)
		return err
	}

	return nil
}

func (s *WebhookServer) processWebhook(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*2500)
	defer cancel()

	httpRequestID := fmt.Sprintf("%v-%v", time.Now().Unix(), atomic.AddUint64(&s.counter, 1))

	ctx, _ = logger.WithField(ctx, "trello_req_id", httpRequestID)
	ctx, _ = logger.WithField(ctx, "method", r.Method)
	ctx, _ = logger.WithField(ctx, "resources", r.RequestURI)
	ctx, _ = logger.WithField(ctx, "from", r.RemoteAddr)

	log := logger.Get(ctx)
	log.Debug("Request received")

	switch r.Method {
	case "GET":
		fmt.Fprintf(rw, "早晨 ~ 😁 \n")
	case "POST":

		body, err := ioutil.ReadAll(r.Body)

		// rdr1 := ioutil.NopCloser(bytes.NewBuffer(body))
		// rdr2 := ioutil.NopCloser(bytes.NewBuffer(body))
		// log.Printf("BODY: %q", rdr1)
		// r.Body = rdr2

		var webhook Webhook

		if err != nil {
			log.WithError(err).Error("Failed to read webhook payload")
			http.Error(rw, "", http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(body, &webhook); err != nil {
			log.WithError(err).Error("Failed to parse webhook payload")
			http.Error(rw, "", http.StatusBadRequest)
			return
		}

		if err = s.onWebhook(ctx, webhook); err != nil {
			log.WithError(err).Error("Failed to process webhook")
			var code int
			switch {
			case lib.IsCanceled(err) || lib.IsDeadline(err):
				code = http.StatusServiceUnavailable
			default:
				code = http.StatusInternalServerError
			}
			http.Error(rw, "", code)
		} else {
			rw.WriteHeader(http.StatusOK)
		}

	default:
		fmt.Fprintf(rw, "%s is not allowed 😛", r.Method)
	}
}
