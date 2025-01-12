package main

import (
	"invoice/data"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var pageTests = []struct {
	name           string
	url            string
	expectedStatus int
	handler        http.HandlerFunc
	sessionData    map[string]any
	expectedHTML   string
}{
	{
		name:           "home",
		url:            "/",
		expectedStatus: http.StatusOK,
		handler:        testApp.HomePage,
	},
	{
		name:           "login",
		url:            "/login",
		expectedStatus: http.StatusOK,
		handler:        testApp.Login,
		expectedHTML:   `<h1 class="mt-5">Login</h1>`,
	},
	{
		name:           "logout",
		url:            "/logout",
		expectedStatus: http.StatusSeeOther,
		handler:        testApp.Logout,
		expectedHTML:   `<h1 class="mt-5">Login</h1>`,
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
	},
}

func Test_Pages(t *testing.T) {
	pathToTemplates = "./templates"

	for _, tc := range pageTests {
		recoder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, tc.url, nil)

		ctx := getCtx(request)
		request = request.WithContext(ctx)

		if len(tc.sessionData) > 0 {
			for key, value := range tc.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		tc.handler.ServeHTTP(recoder, request)

		if status := recoder.Code; status != tc.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
		}

		if len(tc.expectedHTML) > 0 {
			html := recoder.Body.String()
			if !strings.Contains(html, tc.expectedHTML) {
				t.Errorf("handler returned unexpected body: got %v want %v", html, tc.expectedHTML)
			}
		}
	}

}
