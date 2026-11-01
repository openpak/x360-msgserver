// Package msgserver serves the Xbox Live message server (Live service 2,
// /msgserver). Only XStringVerify is served so far. Moved out of x360-xuacs.
package msgserver

import (
	"encoding/binary"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// PathVetString is XStringVerify (XRL 0x801). The dashboard's Edit Profile vets
// Motto and Bio text with it before the profile-setting write; a failed request
// shows "Can't change your motto".
const PathVetString = "/msgserver/vetstring2.ashx"

type Handler struct{ Log *slog.Logger }

func (h *Handler) log() *slog.Logger {
	if h.Log != nil {
		return h.Log
	}
	return slog.Default()
}

// ServeHTTP answers vetString: request LE u32 title, u32 flags, u16 localeLen
// (<=16), u16 count (<=10), locale bytes, count x {u16 len, UTF-8 bytes}.
// Reply: u16 count, count x u32 HRESULT. Every string is accepted (S_OK);
// ponytail: no word filter yet, add one here if the network needs it.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/healthz" {
		return
	}
	if strings.ToLower(r.URL.Path) != PathVetString {
		h.log().Info("unhandled", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, 64<<10))
	if len(body) < 12 {
		http.Error(w, "short request", http.StatusBadRequest)
		return
	}
	n := binary.LittleEndian.Uint16(body[10:])
	if n > 10 {
		http.Error(w, "too many strings", http.StatusBadRequest)
		return
	}
	out := binary.LittleEndian.AppendUint16(nil, n)
	out = append(out, make([]byte, 4*int(n))...)
	h.log().Info("vetstring", "account", r.Header.Get("X-X360-Account"), "strings", n, "body_len", len(body))
	w.Header().Set("Content-Type", "xon/5")
	w.Header().Set("Content-Length", strconv.Itoa(len(out)))
	w.Write(out)
}
