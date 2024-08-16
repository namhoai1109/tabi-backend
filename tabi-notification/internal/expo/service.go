package expo

import (
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

type Expo interface {
	SendNotification(token, title, body string) error
	SendNotifications(tokens []string, title, body string) error
}

func New() Expo {
	return &service{
		client: expo.NewPushClient(nil),
	}
}

type service struct {
	client *expo.PushClient
}
