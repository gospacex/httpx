package hertz

import (
	"context"
	"testing"
)

func TestHertzServer_StartStop(t *testing.T) {
	srv := NewServer()

	if srv.IsRunning() {
		t.Error("Server should not be running initially")
	}
}

func TestHertzServer_GracefulShutdown(t *testing.T) {
	srv := NewServer()

	srv.GracefulShutdown(context.Background())

	if srv.IsRunning() {
		t.Error("Server should be stopped after GracefulShutdown")
	}
}