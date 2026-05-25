package main

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gospacex/httpx"
	_ "github.com/gospacex/httpx/adapter/gin"
)

const ginPort = 18891

var ginBaseURL string

func ginDoReq(method, path, body string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, ginBaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func startGinServer(configPath string, port int) {
	ginBaseURL = "http://localhost:" + strconv.Itoa(port)
	app, err := httpx.New(configPath)
	if err != nil {
		panic(err)
	}
	setupRoutes(app)
	go func() {
		app.RunOnAddr(":" + strconv.Itoa(port))
	}()
	time.Sleep(800 * time.Millisecond)
}

func TestGinAdapter(t *testing.T) {
	startGinServer("./config.yaml", ginPort)

	resp, err := ginDoReq("GET", "/", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Root: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/health", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Health: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/hello?name=test", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Hello: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/time", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Time: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/api/v1/info", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("APIV1Info: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/api/v1/timestamp", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("APIV1Timestamp: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/admin/dashboard", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("AdminUnauthorized: expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/admin/dashboard?token=admin123", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("AdminAuthorized: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/users", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("UsersList: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	body := `{"name":"test","age":25}`
	resp, err = ginDoReq("POST", "/users", body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("UsersCreate: expected 201, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("GET", "/articles", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("ArticlesList: expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	articleBody := `{"title":"test","content":"hello","author":"tester"}`
	resp, err = ginDoReq("POST", "/articles", articleBody)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("ArticlesCreate: expected 201, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("PUT", "/users/1", `{"name":"updated","age":30}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		t.Fatalf("UserUpdate: expected 200 or 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("DELETE", "/users/1", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		t.Fatalf("UserDelete: expected 200 or 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp, err = ginDoReq("PATCH", "/users/1", `{"age":35}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		t.Fatalf("PATCHEndpoint: expected 200 or 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}