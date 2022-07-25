package sender

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"strconv"
	"time"
)

type Email struct {
	c config
}
type config struct {
	host     string
	port     int
	username string
	password string
	ssl      bool
	tls      bool
}

func (e Email) Send(subject string, message []byte, from string, to []string) error {
	server := mail.NewSMTPClient()
	// SMTP Server
	server.Host = e.c.host
	server.Port = e.c.port
	server.Username = e.c.username
	server.Password = e.c.password
	switch {
	case e.c.ssl && e.c.tls:
		server.Encryption = mail.EncryptionSSLTLS
	case e.c.ssl:
		server.Encryption = mail.EncryptionSSL
	case e.c.tls:
		server.Encryption = mail.EncryptionTLS
	default:
		server.Encryption = mail.EncryptionNone
	}
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("%v <%v>", "SHIO", from))
	email.AddTo(to...).
		SetSubject(subject).
		SetBody(mail.TextHTML, string(message))
	return email.Send(smtpClient)
}
func NewEmailSender() Email {
	c, err := loadEnv()
	if err != nil {
		log.Fatalln(err)
	}
	return Email{c}
}
func loadEnv() (cfg config, err error) {
	const (
		HOST = iota
		PORT
		SSL
		TLS
		USERNAME
		PASSWORD
	)
	need := [6]string{"SMTP_HOST", "SMTP_PORT", "SMTP_SSL", "SMTP_TLS", "SMTP_USERNAME", "SMTP_PASSWORD"}
	val := [6]string{}
	for i, el := range need {
		if v, exists := os.LookupEnv(el); !exists {
			return cfg, fmt.Errorf("environment %v must be not empty", need[i])
		} else {
			val[i] = v
		}
	}
	portI, err := strconv.Atoi(val[PORT])
	if err != nil {
		return cfg, errors.New("port must be integer")
	}
	sslB, err := strconv.ParseBool(val[SSL])
	if err != nil {
		return cfg, fmt.Errorf("%v must be bool", need[SSL])
	}
	tlsB, err := strconv.ParseBool(val[TLS])
	if err != nil {
		return cfg, fmt.Errorf("%v must be bool", need[TLS])
	}
	return config{
		host:     val[HOST],
		port:     portI,
		ssl:      sslB,
		tls:      tlsB,
		username: val[USERNAME],
		password: val[PASSWORD],
	}, nil
}
