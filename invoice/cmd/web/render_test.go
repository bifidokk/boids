package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig_AddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	templateData := testApp.AddDefaultData(&TemplateData{}, req)

	if templateData.Flash != "flash" {
		t.Errorf("TemplateData was incorrect, got: %s, want: %s.", templateData.Flash, "flash")
	}

	if templateData.Warning != "warning" {
		t.Errorf("TemplateData was incorrect, got: %s, want: %s.", templateData.Warning, "warning")
	}

	if templateData.Error != "error" {
		t.Errorf("TemplateData was incorrect, got: %s, want: %s.", templateData.Error, "error")
	}
}

func TestConfig_IsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	auth := testApp.IsAuthenticated(req)

	if auth {
		t.Errorf("IsAuthenticated was incorrect, got: %v, want: %t.", auth, false)
	}

	testApp.Session.Put(ctx, "userID", 1)

	auth = testApp.IsAuthenticated(req)

	if !auth {
		t.Errorf("IsAuthenticated was incorrect, got: %v, want: %t.", auth, true)
	}
}

func TestConfig_render(t *testing.T) {
	pathToTemplates = "./templates"

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.render(recorder, req, "home.page.gohtml", &TemplateData{})

	if recorder.Code != http.StatusOK {
		t.Errorf("render was incorrect, got: %d, want: %d.", recorder.Code, http.StatusOK)
	}
}
