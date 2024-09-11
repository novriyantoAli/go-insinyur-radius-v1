package mbroker

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type messageHandler struct {
	ucase domain.MessageUsecase
	rc    *redis.Client
}

// NewHandler ...
func NewHandler(rc *redis.Client, uc domain.MessageUsecase) {
	handler := &messageHandler{ucase: uc, rc: rc}

	go handler.SubscribeUser()
}

func (hn *messageHandler) SubscribeUser() error {
	subscriber := hn.rc.Subscribe(context.Background(), "user-data")
	user := domain.User{}

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			logrus.Error(err)
			return err
		}

		if err := json.Unmarshal([]byte(msg.Payload), &user); err != nil {
			logrus.Error(err)
			return err
		}

		if msg.Channel == "user-data" {
			err = hn.ucase.ReadUserMessage(context.Background(), &user)

		}

	}
}
