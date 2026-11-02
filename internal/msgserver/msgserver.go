// Package msgserver serves the Xbox Live message server (Live service 2,
// /msgserver). Only XStringVerify is served so far. Moved out of x360-xuacs.
package msgserver

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// PathTitles is the XamGetCachedTitleName lookup (XRL 0x800); PathStrings is the string lookup (0x802).
const (
	PathTitles  = "/msgserver/msgserver.ashx"
	PathStrings = "/msgserver/getstring2.ashx"
)

// LoadTitles reads a JSON object of title id (hex) -> name. An empty path gives an empty map.
// ponytail: static file; the website catalog is the source later
func LoadTitles(path string) (map[uint32]string, error) {
	m := map[uint32]string{}
	if path == "" {
		return m, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]string
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	for k, v := range raw {
		var id uint32
		if _, err := fmt.Sscanf(k, "%x", &id); err != nil {
			return nil, fmt.Errorf("title id %q: %w", k, err)
		}
		m[id] = v
	}
	return m, nil
}

// PathVetString is XStringVerify (XRL 0x801). The dashboard's Edit Profile vets
// Motto and Bio text with it before the profile-setting write; a failed request
// shows "Can't change your motto".
const PathVetString = "/msgserver/vetstring2.ashx"

type Handler struct {
	Log    *slog.Logger
	Titles map[uint32]string
}

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
	switch strings.ToLower(r.URL.Path) {
	case PathTitles:
		h.names(w, r, false)
		return
	case PathStrings:
		h.names(w, r, true)
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

// names answers 0x800 (request u16 1, u16 key, u16 n <=100, n x u32 title id) and 0x802 (u32 a0, u16 n <=100,
// u16 locale length <=16, locale, n x u32 string id). Reply u16 n, n x {u16 len, UTF-8 name}; an unknown id
// gets len 0. No string table is known, so 0x802 answers every id empty.
func (h *Handler) names(w http.ResponseWriter, r *http.Request, strs bool) {
	body, _ := io.ReadAll(io.LimitReader(r.Body, 64<<10))
	var n, off int
	if strs {
		if len(body) < 8 {
			http.Error(w, "short request", http.StatusBadRequest)
			return
		}
		n = int(binary.LittleEndian.Uint16(body[4:]))
		off = 8 + int(binary.LittleEndian.Uint16(body[6:]))
		if binary.LittleEndian.Uint16(body[6:]) > 16 {
			n = 101
		}
	} else {
		if len(body) < 6 {
			http.Error(w, "short request", http.StatusBadRequest)
			return
		}
		n = int(binary.LittleEndian.Uint16(body[4:]))
		off = 6
	}
	if n > 100 || len(body) < off+4*n {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	out := binary.LittleEndian.AppendUint16(nil, uint16(n))
	for i := 0; i < n; i++ {
		name := ""
		if !strs {
			name = h.Titles[binary.LittleEndian.Uint32(body[off+4*i:])]
		}
		out = binary.LittleEndian.AppendUint16(out, uint16(len(name)))
		out = append(out, name...)
	}
	h.log().Info("names", "strings", strs, "account", r.Header.Get("X-X360-Account"), "n", n)
	w.Header().Set("Content-Type", "xon/5")
	w.Header().Set("Content-Length", strconv.Itoa(len(out)))
	w.Write(out)
}
