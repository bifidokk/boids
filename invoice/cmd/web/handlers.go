package main

import (
	"fmt"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"html/template"
	"invoice/data"
	"net/http"
	"strconv"
	"time"
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

	validPassword, err := app.Models.User.PasswordMatches(password)

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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (app *Config) PlansList(w http.ResponseWriter, r *http.Request) {
	if !app.Session.Exists(r.Context(), "userID") {
		app.Session.Put(r.Context(), "error", "Invalid session")
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	dataMap := map[string]any{
		"plans": plans,
	}

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func (app *Config) Subscribe(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	plan, err := app.Models.Plan.GetOne(id)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	user, ok := app.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		app.ErrorLog.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()

		invoice, err := app.getInvoice(user, plan)
		if err != nil {
			app.ErrorChan <- err
			return
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Invoice subscription",
			Template: "invoice",
			Data:     invoice,
		}

		app.sendEmail(msg)
	}()

	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()

		pdf := app.generateManual(user, plan)
		filename := fmt.Sprintf("./tmp/%d_manual.pdf", user.ID)

		err := pdf.OutputFileAndClose(filename)
		if err != nil {
			app.ErrorChan <- err
			return
		}

		msg := Message{
			To:      user.Email,
			Subject: "You manual",
			Data:    "You manual is activated",
			AttachmentMap: map[string]string{
				"manual.pdf": filename,
			},
		}

		app.sendEmail(msg)
	}()

	err = app.Models.Plan.SubscribeUserToPlan(user, *plan)
	if err != nil {
		app.ErrorChan <- err
		http.Redirect(w, r, "/members/plans", http.StatusFound)
		return
	}

	updatedUser, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		app.ErrorChan <- err
		http.Redirect(w, r, "/members/plans", http.StatusFound)
		return
	}

	app.Session.Put(r.Context(), "user", updatedUser)

	app.Session.Put(r.Context(), "flash", "Subscription activated")
	http.Redirect(w, r, "/members/plans", http.StatusFound)
}

func (app *Config) getInvoice(user data.User, plan *data.Plan) (string, error) {
	return plan.PlanAmountFormatted, nil
}

func (app *Config) generateManual(user data.User, plan *data.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	time.Sleep(5 * time.Second)

	templateId := importer.ImportPage(pdf, "./pdf/manual.pdf", 1, "/MediaBox")
	pdf.AddPage()

	importer.UseImportedTemplate(pdf, templateId, 0, 0, 215.9, 0)

	pdf.SetX(75)
	pdf.SetY(150)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", user.FirstName, user.LastName), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", user.Email), "", "C", false)

	return pdf
}
