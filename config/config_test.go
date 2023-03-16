package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/0xERR0R/mailcatcher/config"
)

var _ = Describe("Config", func() {
	Describe("Get configuration from ENV", func() {
		BeforeEach(func() {
			DeferCleanup(os.Setenv, "MC_PORT", os.Getenv("MC_PORT"))

		})
		It("Should create valide configuraiotn", func() {
			os.Setenv("MC_PORT", "112")
			os.Setenv("MC_HOST", "domain.tld")
			os.Setenv("MC_REDIRECT_TO", "user1@example.com")
			os.Setenv("MC_SENDER_MAIL", "user2@example.com")
			os.Setenv("MC_SMTP_HOST", "imap.gmail.com")
			os.Setenv("MC_SMTP_PORT", "25")
			os.Setenv("MC_SMTP_USER", "user")
			os.Setenv("MC_SMTP_PASSWORD", "pw")

			cfg, err := config.GetConfiguration()

			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.Port).Should(Equal(112))
			Expect(cfg.Host).Should(Equal("domain.tld"))
			Expect(cfg.RedirectTo).Should(Equal("user1@example.com"))
			Expect(cfg.SenderMail).Should(Equal("user2@example.com"))
			Expect(cfg.SMTPHost).Should(Equal("imap.gmail.com"))
			Expect(cfg.SMTPPort).Should(Equal(25))
			Expect(cfg.SMTPUser).Should(Equal("user"))
			Expect(cfg.SMTPPassword).Should(Equal("pw"))
		})
	})
})
