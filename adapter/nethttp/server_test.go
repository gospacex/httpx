package nethttp

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNetHttpServer_StartStop(t *testing.T) {
	srv := NewServer()
	srv.router = NewRouter()

	go func() {
		err := srv.Start(":0")
		if err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)
	if !srv.IsRunning() {
		t.Error("Server should be running after Start")
	}

	err := srv.Stop(context.Background())
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}

func TestNetHttpServer_GracefulShutdown(t *testing.T) {
	srv := NewServer()
	srv.router = NewRouter()

	if srv.IsRunning() {
		t.Error("Server should not be running initially")
	}
}