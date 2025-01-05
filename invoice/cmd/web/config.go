package main

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"invoice/data"
	"log"
	"sync"
)

type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
	Mailer   Mail
}

func (app *Config) createMailer() Mail {
	errorChannel := make(chan error)
	mailerChannel := make(chan Message, 100)
	doneChannel := make(chan bool)

	m := Mail{
		Domain:        "localhost",
		Host:          "localhost",
		Port:          1025,
		Encryption:    "none",
		FromAddress:   "info@example.com",
		FromName:      "info",
		Wait:          app.Wait,
		ErrorChannel:  errorChannel,
		DoneChannel:   doneChannel,
		MailerChannel: mailerChannel,
	}

	return m
}
