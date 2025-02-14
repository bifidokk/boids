package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"invoice/data"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort = 80
)

func main() {
	db := initDb()
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	session := initSession()
	wg := &sync.WaitGroup{}
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := Config{
		Session:       session,
		DB:            db,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Wait:          wg,
		Models:        data.New(db),
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}

	app.Mailer = app.createMailer()
	go app.listenForEmail()
	go app.listenForErrors()
	go app.listenForShutdown()

	app.serve()
}

func (app *Config) serve() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Printf("Starting server on port %d", webPort)
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func initDb() *sql.DB {
	conn := connectToDB()

	if conn == nil {
		panic("failed to connect to database")
	}

	return conn
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	connection, err := openDB(dsn)

	if err != nil {
		log.Fatal(err)
	}

	return connection
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSession() *scs.SessionManager {
	gob.Register(data.User{})
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}

	return pool
}

func (app *Config) listenForErrors() {
	for {
		select {
		case err := <-app.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.ErrorChanDone:
			return
		}
	}
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	app.InfoLog.Println("Shutting down...")

	app.Wait.Wait()

	app.Mailer.DoneChannel <- true
	app.ErrorChanDone <- true

	close(app.Mailer.DoneChannel)
	close(app.Mailer.ErrorChannel)
	close(app.Mailer.MailerChannel)
	close(app.ErrorChan)
	close(app.ErrorChanDone)
}
