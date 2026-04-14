package email

type EmailSender interface {
	Send(to string, code string) error
}
