package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"invoice/data"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

var testApp Config

func TestMain(m *testing.M) {
	gob.Register(data.User{})

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	testApp = Config{
		Session:       session,
		DB:            nil,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Wait:          &sync.WaitGroup{},
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}

	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	testApp.Mailer = Mail{
		Wait:          testApp.Wait,
		ErrorChannel:  errorChan,
		MailerChannel: mailerChan,
		DoneChannel:   mailerDoneChan,
	}

	go func() {
		select {
		case <-testApp.Mailer.MailerChannel:
		case <-testApp.Mailer.ErrorChannel:
		case <-testApp.Mailer.DoneChannel:
			return
		}
	}()

	go func() {
		for {
			select {
			case err := <-errorChan:
				testApp.ErrorLog.Println(err)
			case <-testApp.ErrorChanDone:
				return
			}
		}
	}()

	os.Exit(m.Run())
}
