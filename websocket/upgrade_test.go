package websocket

import (
	"net/http"
	"testing"
)

func TestUpgraderStruct(t *testing.T) {
	u := &Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if u.ReadBufferSize != 1024 {
		t.Errorf("ReadBufferSize = %d, want 1024", u.ReadBufferSize)
	}
	if u.WriteBufferSize != 1024 {
		t.Errorf("WriteBufferSize = %d, want 1024", u.WriteBufferSize)
	}
	if u.CheckOrigin == nil {
		t.Error("CheckOrigin is nil")
	}
}

func TestDefaultUpgrader(t *testing.T) {
	if DefaultUpgrader == nil {
		t.Fatal("DefaultUpgrader is nil")
	}
	if DefaultUpgrader.ReadBufferSize != 1024 {
		t.Errorf("DefaultUpgrader.ReadBufferSize = %d, want 1024", DefaultUpgrader.ReadBufferSize)
	}
	if DefaultUpgrader.WriteBufferSize != 1024 {
		t.Errorf("DefaultUpgrader.WriteBufferSize = %d, want 1024", DefaultUpgrader.WriteBufferSize)
	}
}