package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/72sevenzy2/grouter/grouter"
	utils "github.com/72sevenzy2/production-backend-template/internal"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

const SESSIONID_EXPIRY_TIME = time.Minute * 10

func HandleUserLogin(u *repo.Users, sessions *repo.Sessions) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		limit_body := io.LimitReader(r.Body, 1024*1024) // 1 MiB

		var LoginReq utils.LoginRequest
		if err := json.NewDecoder(limit_body).Decode(&LoginReq); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		user, ok := u.GetUser(LoginReq.ID)
		if !ok {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusFailedDependency,
				Title:      fmt.Sprintf("status %d", http.StatusFailedDependency),
				Heading:    fmt.Sprintf("status %d", http.StatusFailedDependency),
				Descripion: "Sorry, your request failed as the user you were looking for was not found.",
			})
			return
		}

		sID, err := generateSessionID()
		if err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		// password comparison
		same, err := utils.VerifyPassword(LoginReq.Password, user.GetHashedPassword())
		if err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}
		if !same {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusUnauthorized,
				Title:      fmt.Sprintf("status %d", http.StatusUnauthorized),
				Heading:    fmt.Sprintf("status %d", http.StatusUnauthorized),
				Descripion: "Sorry, the password you entered for the account was wrong.",
			})
			return
		}

		sessions.Create(LoginReq.ID, sID, time.Now().Add(SESSIONID_EXPIRY_TIME))

		// some confirmation responses
		returning := struct {
			Username string `json:"username"`
			Id       int    `json:"id"`
		}{
			Username: user.GetName(),
			Id:       user.GetID(),
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sID,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			MaxAge:   int(time.Duration(SESSIONID_EXPIRY_TIME).Seconds()),
		})

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(returning); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
