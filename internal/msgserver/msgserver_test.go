package msgserver

import (
	"bytes"
	"encoding/hex"
	"net/http/httptest"
	"testing"
)

func TestVetStringAcceptsAll(t *testing.T) {
	// title FFFE07D1, flags 0, locale "en-US", two strings "hi" and "n".
	body, _ := hex.DecodeString("d107feff0000000005000200656e2d5553020068690100" + "6e")
	rec := httptest.NewRecorder()
	(&Handler{}).ServeHTTP(rec, httptest.NewRequest("POST", PathVetString, bytes.NewReader(body)))
	if got := hex.EncodeToString(rec.Body.Bytes()); rec.Code != 200 || got != "02000000000000000000" {
		t.Fatalf("reply %d %s", rec.Code, got)
	}
	rec = httptest.NewRecorder()
	(&Handler{}).ServeHTTP(rec, httptest.NewRequest("POST", "/msgserver/other.ashx", bytes.NewReader(body)))
	if rec.Code != 404 {
		t.Fatal("other path served")
	}
}
