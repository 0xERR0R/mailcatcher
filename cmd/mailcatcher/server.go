package main

import (
	"io"
	"log"
	"net/smtp"
	"strconv"
	"strings"

	"fmt"

	gosmtp "github.com/emersion/go-smtp"
	"github.com/veqryn/go-email/email"
)

var config *Configuration

type Backend struct{}

func (bkd *Backend) NewSession(c *gosmtp.Conn) (gosmtp.Session, error) {
	return &Session{
		to: make([]string, 0),
	}, nil
}

type Session struct {
	from string
	to   []string
}

func (s *Session) Mail(from string, opts *gosmtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *gosmtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	log.Printf("New message from '%s' to '%s' received", s.from, s.to)
	if isRecipientValid(s.to) {
		if msg, err := email.ParseMessage(r); err != nil {
			log.Fatal("error", err)
			return err
		} else {
			msg.Header.SetSubject(fmt.Sprintf("[MAILCATCHER] %s", msg.Header.Subject()))
			msg.Header.SetTo(fmt.Sprintf("\"%s\" <%s>", msg.Header.To()[0], config.MC_REDIRECT_TO))
			msg.Header.SetFrom(fmt.Sprintf("\"%s\" <%s>", "MAILCATCHER", config.MC_SENDER_MAIL))

			sendMail(msg)

			if err != nil {
				log.Printf("smtp error: %s", err)
			}
		}
	} else {
		log.Print("ignoring message")
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func isRecipientValid(recipients []string) bool {
	for _, recipient := range recipients {
		if strings.HasSuffix(recipient, config.MC_HOST) {
			return true
		}
	}
	return false
}

func sendMail(msg *email.Message) {
	if err := msg.Save(); err != nil {
		log.Printf("can't save message: %s", err)
		return
	}
	b, err := msg.Bytes()
	if err != nil {
		log.Printf("can't convert message: %s", err)
		return
	}

	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.MC_SMTP_HOST, config.MC_SMTP_PORT),
		smtp.PlainAuth("", config.MC_SMTP_USER, config.MC_SMTP_PASSWORD, config.MC_SMTP_HOST),
		config.MC_SENDER_MAIL, []string{config.MC_REDIRECT_TO}, b)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func NewServer(configuration *Configuration) error {
	config = configuration
	be := &Backend{}

	s := gosmtp.NewServer(be)

	s.Addr = ":" + strconv.Itoa(config.MC_PORT)
	s.Domain = config.MC_HOST
	s.MaxMessageBytes = 1024 * 1024 * 20
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	return s.ListenAndServe()
}
