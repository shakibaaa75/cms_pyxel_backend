package handlers

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"time"

	"cms-backend/config"
	"cms-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONErr(w, "invalid request body", http.StatusBadRequest)
		return
	}

	usernameMatch := subtle.ConstantTimeCompare([]byte(req.Username), []byte(config.ADMIN_USERNAME)) == 1
	bcryptErr := bcrypt.CompareHashAndPassword(config.AdminPasswordHash, []byte(req.Password))

	if !usernameMatch || bcryptErr != nil {
		time.Sleep(200 * time.Millisecond)
		utils.JSONErr(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := issueAdminJWT()
	if err != nil {
		utils.JSONErr(w, "internal server error", http.StatusInternalServerError)
		return
	}
	utils.JSON200(w, map[string]string{"token": token, "message": "login successful"})
}

func issueAdminJWT() (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   "admin",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		Issuer:    config.COMPANY_NAME,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.JWT_SECRET)
}
