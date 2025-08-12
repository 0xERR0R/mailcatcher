package main

import (
	"fmt"
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type Configuration struct {
	MC_PORT          int    `validate:"required,gte=0,lte=65535"`
	MC_HOST          string `validate:"required,hostname"`
	MC_REDIRECT_TO   string `validate:"required,email"`
	MC_SENDER_MAIL   string `validate:"required,email"`
	MC_SMTP_HOST     string `validate:"required,hostname"`
	MC_SMTP_PORT     int    `validate:"required,gte=0,lte=65535"`
	MC_SMTP_USER     string `validate:"omitempty"`
	MC_SMTP_PASSWORD string `validate:"omitempty"`
}

func (c *Configuration) Validate() error {

	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(c)
	if err != nil {
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
		c.MC_PORT, c.MC_HOST, c.MC_REDIRECT_TO, c.MC_SENDER_MAIL, c.MC_SMTP_HOST, c.MC_SMTP_PORT, c.MC_SMTP_USER)
}
