package mock

import (
	"control_plane/internal/service/email"
	"log/slog"
)

type DevEmailSender struct {
	log *slog.Logger
}

func NewEmailSenderMock(log *slog.Logger) email.EmailSender {
	return &DevEmailSender{log: log}
}

func (s *DevEmailSender) Send(to string, code string) error {
	s.log.Info("DEV EMAIL",
		"to", to,
		"code", code,
	)
	return nil
}
