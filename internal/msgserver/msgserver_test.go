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

func TestTitleNames(t *testing.T) {
	h := &Handler{Titles: map[uint32]string{0xFFFE07D1: "Xbox 360 Dashboard"}}
	// u16 1, key 0, n 2: FFFE07D1 and an unknown id.
	body, _ := hex.DecodeString("010000000200" + "d107feff" + "01000000")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", PathTitles, bytes.NewReader(body)))
	want := "0200" + "1200" + hex.EncodeToString([]byte("Xbox 360 Dashboard")) + "0000"
	if got := hex.EncodeToString(rec.Body.Bytes()); rec.Code != 200 || got != want {
		t.Fatalf("reply %d %s", rec.Code, got)
	}
}

func TestStringLookupEmpty(t *testing.T) {
	// a0 0, n 2, locale "en-US", ids 1 and 2.
	body, _ := hex.DecodeString("00000000" + "0200" + "0500" + "656e2d5553" + "01000000" + "02000000")
	rec := httptest.NewRecorder()
	(&Handler{}).ServeHTTP(rec, httptest.NewRequest("POST", PathStrings, bytes.NewReader(body)))
	if got := hex.EncodeToString(rec.Body.Bytes()); rec.Code != 200 || got != "020000000000" {
		t.Fatalf("reply %d %s", rec.Code, got)
	}
}
