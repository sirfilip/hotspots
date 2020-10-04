package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"necsam/config"
	"necsam/models"
)

var (
	ActivationCodeTemplateEng = template.Must(template.New("ActivationCodeEng").Parse(`
		Dear {{ .Username }}

		You have registered at {{ .App }}.

		Please ignore this email if that wasn't you.

		To activate your account please follow {{ .ActivationURL }}
	`))
)

type ActivationCodeSender interface {
	SendActivationCode(user models.User, activationCode string) error
}

type Mailer struct {
	EmailHost         string
	EmailHostUser     string
	EmailHostPassword string
	EmailPort         int64
	EmailUseTls       bool
}

func (mailer Mailer) SendActivationCode(user models.User, activationCode string) error {
	buf := &bytes.Buffer{}

	err := ActivationCodeTemplateEng.ExecuteTemplate(buf, "ActivationCodeEng", map[string]string{
		"Username":      user.Username,
		"App":           config.Get("app_name"),
		"ActivationURL": config.Get("api_url") + "/activate/" + activationCode,
	})
	err = mailer.send("Activation Code", []string{user.Email}, buf.Bytes())

	return err
}

func (mailer Mailer) send(subject string, to []string, msg []byte) error {
	var err error
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("From: %s\r\n", config.Get("email_host_user")))
	b.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ",")))
	b.WriteString(fmt.Sprintf("Subject: %s\r\n\r\n", subject))
	b.WriteString(string(msg))

	emailHost := config.Get("email_host")
	emailPort := config.Get("email_port")
	emailServerAddr := fmt.Sprintf("%s:%s", emailHost, emailPort)
	emailHostUser := config.Get("email_host_user")
	emailHostPassword := config.Get("email_host_password")
	emailUseTls := config.GetBool("email_use_tls")
	emailTlsInsecure := config.GetBool("email_tls_skip_verify")

	c, err := smtp.Dial(emailServerAddr)
	if err != nil {
		return err
	}

	if emailUseTls {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: emailTlsInsecure,
			ServerName:         emailHost,
		}
		c.StartTLS(tlsconfig)
		auth := smtp.PlainAuth("", emailHostUser, emailHostPassword, emailHost)
		if err := c.Auth(auth); err != nil {
			return err
		}
	}

	if err := c.Mail(emailHostUser); err != nil {
		return err
	}

	if err := c.Rcpt(strings.Join(to, ",")); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write(b.Bytes()); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	c.Quit()

	return nil
}
