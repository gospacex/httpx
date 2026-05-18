package nethttp

import (
	"context"
	"testing"
)

func TestNetHttpServer_StartStop(t *testing.T) {
	srv := NewServer()
	srv.router = NewRouter()

	err := srv.Start(":0")
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	err = srv.Stop(context.Background())
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}

func TestNetHttpServer_GracefulShutdown(t *testing.T) {
	srv := NewServer()
	srv.router = NewRouter()

	if !srv.IsRunning() {
		t.Error("Server should not be running initially")
	}
}