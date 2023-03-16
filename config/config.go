package config

import (
	"fmt"
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type Configuration struct {
	Port         int    `env:"MC_PORT" validate:"required,gte=0,lte=65535"`
	Host         string `env:"MC_HOST" validate:"required,hostname"`
	RedirectTo   string `env:"MC_REDIRECT_TO" validate:"required,email"`
	SenderMail   string `env:"MC_SENDER_MAIL" validate:"required,email"`
	SMTPHost     string `env:"MC_SMTP_HOST" validate:"required,hostname"`
	SMTPPort     int    `env:"MC_SMTP_PORT" validate:"required,gte=0,lte=65535"`
	SMTPUser     string `env:"MC_SMTP_USER" validate:"required"`
	SMTPPassword string `env:"MC_SMTP_PASSWORD" validate:"required"`
}

func (c *Configuration) Validate() error {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(c)
	if err != nil {
		//nolint:errorlint
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			log.Printf("configuration error: %s", e.Translate(trans))
		}
	}

	return err
}

func (c Configuration) String() string {
	return fmt.Sprintf(`
	MC_PORT:          %d
	MC_HOST:          %s
	MC_REDIRECT_TO:   %s
	MC_SENDER_MAIL:   %s
	MC_SMTP_HOST:     %s
	MC_SMTP_PORT:     %d
	MC_SMTP_USER:     %s`,
		c.Port, c.Host, c.RedirectTo, c.SenderMail, c.SMTPHost, c.SMTPPort, c.SMTPUser)
}
