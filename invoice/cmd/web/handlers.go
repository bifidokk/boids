package main

import (
	"fmt"
	"html/template"
	"invoice/data"
	"net/http"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	validPassword, err := user.PasswordMatches(password)

	if err != nil || !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed log in attempt",
			Data:    "Failed log in attempt",
		}

		app.sendEmail(msg)

		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	app.Session.Put(r.Context(), "flash", "Logged in")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	app.Session.Put(r.Context(), "flash", "Logged out")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
		http.Redirect(w, r, "/register", http.StatusFound)

		return
	}

	user := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}

	_, err = user.Insert(user)

	if err != nil {
		app.ErrorLog.Println(err)
		http.Redirect(w, r, "/register", http.StatusFound)

		return
	}

	url := fmt.Sprintf("http://localhost/activate?email=%s", user.Email)
	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       user.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.sendEmail(msg)
	app.Session.Put(r.Context(), "flash", "Activate your account")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	testURL := fmt.Sprintf("http://localhost%s", url)
	ok := VerifyToken(testURL)

	if !ok {
		app.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	user, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid email")
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	user.Active = 1
	err = user.Update()

	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to update user")
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	app.Session.Put(r.Context(), "flash", "Account activated")
	http.Redirect(w, r, "/login", http.StatusFound)
}
