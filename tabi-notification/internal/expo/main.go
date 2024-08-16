package expo

import (
	"fmt"

	"github.com/namhoai1109/tabi/core/server"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func (s *service) SendNotification(token, title, body string) error {
	// To check the token is valid
	pushToken, err := expo.NewExponentPushToken(token)
	if err != nil {
		return server.NewHTTPInternalError(fmt.Sprintf("Invalid token: %s", token))
	}

	// Publish message
	response, err := s.client.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{pushToken},
			Body:     body,
			Sound:    "default",
			Title:    title,
			Priority: expo.DefaultPriority,
		},
	)

	// Check errors
	if err != nil {
		return server.NewHTTPInternalError(fmt.Sprintf("Failed to send notification with err %s", err.Error()))
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		return server.NewHTTPInternalError(fmt.Sprintf("Failed to send notification with err %s", response.ValidateResponse().Error()))
	}

	return nil
}

func (s *service) SendNotifications(tokens []string, title, body string) error {
	pushTokens := make([]expo.ExponentPushToken, 0)

	for _, token := range tokens {
		// To check the token is valid
		pushToken, err := expo.NewExponentPushToken(token)
		if err != nil {
			return server.NewHTTPInternalError(fmt.Sprintf("Invalid token: %s", token))
		}

		pushTokens = append(pushTokens, pushToken)
	}

	// Publish message
	response, err := s.client.Publish(
		&expo.PushMessage{
			To:       pushTokens,
			Body:     body,
			Sound:    "default",
			Title:    title,
			Priority: expo.DefaultPriority,
		},
	)

	// Check errors
	if err != nil {
		return server.NewHTTPInternalError(fmt.Sprintf("Failed to send notification with err %s", err.Error()))
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		return server.NewHTTPInternalError(fmt.Sprintf("Failed to send notification with err %s", response.ValidateResponse().Error()))
	}

	return nil
}
