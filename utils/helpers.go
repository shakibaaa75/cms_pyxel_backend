package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters   = map[string]*ipLimiter{}
	limitersMu sync.Mutex
)

func RealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.SplitN(ip, ",", 2)[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	addr := r.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i >= 0 {
		return addr[:i]
	}
	return addr
}

func GetLimiter(ip string) *rate.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()
	entry, ok := limiters[ip]
	if !ok {
		entry = &ipLimiter{limiter: rate.NewLimiter(rate.Every(time.Minute), 60)}
		limiters[ip] = entry
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

func PruneOldLimiters() {
	for range time.Tick(5 * time.Minute) {
		limitersMu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for ip, e := range limiters {
			if e.lastSeen.Before(cutoff) {
				delete(limiters, ip)
			}
		}
		limitersMu.Unlock()
	}
}

func GenerateSecureID() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	result := make([]byte, 8)
	for i, v := range b {
		result[i] = chars[int(v)%len(chars)]
	}
	return string(result)
}

func GenerateAccessCode() string {
	return "PXL-" + GenerateSecureID()[:6]
}

func GenerateMagicToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}

func JSON200(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func JSONErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func ExtractPathSegment(path, prefix string, n int) string {
	trimmed := strings.TrimPrefix(path, prefix)
	parts := strings.Split(trimmed, "/")
	if n >= len(parts) {
		return ""
	}
	return parts[n]
}

func ParseObjectID(w http.ResponseWriter, raw string) (bson.ObjectID, bool) {
	id, err := bson.ObjectIDFromHex(raw)
	if err != nil {
		JSONErr(w, "invalid id", http.StatusBadRequest)
		return bson.ObjectID{}, false
	}
	return id, true
}

func RequireSiteID(w http.ResponseWriter, r *http.Request) (string, bool) {
	siteID := strings.TrimSpace(r.URL.Query().Get("siteId"))
	if siteID == "" {
		JSONErr(w, "siteId required", http.StatusBadRequest)
		return "", false
	}
	return siteID, true
}
