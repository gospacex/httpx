package main

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gospacex/httpx"
	_ "github.com/gospacex/httpx/adapter/hertz"
)

const hertzPort = 8888

var hertzBaseURL string

func hertzDoReq(method, path, body string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, hertzBaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func startHertzServer(configPath string) {
	hertzBaseURL = "http://localhost:" + strconv.Itoa(hertzPort)
	app, err := httpx.New(configPath)
	if err != nil {
		panic(err)
	}
	setupRoutes(app)
	go func() {
		app.Run()
	}()
	time.Sleep(800 * time.Millisecond)
}

func TestHertzAdapter(t *testing.T) {
	startHertzServer("./config_hertz.yaml")

	resp, err := hertzDoReq("GET", "/", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Root: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/health", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Health: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/hello?name=hertz", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Hello: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/time", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Time: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/api/v1/info", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("APIV1Info: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/api/v1/timestamp", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("APIV1Timestamp: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/admin/dashboard", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("AdminUnauthorized: expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/admin/dashboard?token=admin123", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("AdminAuthorized: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/users", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("UsersList: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	body := `{"name":"test","age":25}`
	resp, err = hertzDoReq("POST", "/users", body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("UsersCreate: expected 201, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("GET", "/articles", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("ArticlesList: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	articleBody := `{"title":"test","content":"hello","author":"tester"}`
	resp, err = hertzDoReq("POST", "/articles", articleBody)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("ArticlesCreate: expected 201, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("PUT", "/users/1", `{"name":"updated","age":30}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		t.Fatalf("UserUpdate: expected 200 or 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("DELETE", "/users/1", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		t.Fatalf("UserDelete: expected 200 or 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = hertzDoReq("PATCH", "/users/1", `{"age":35}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		t.Fatalf("PATCHEndpoint: expected 200 or 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}