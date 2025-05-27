package handlers

import (
	"context"
	"encoding/json"
	"github/eventApp/internal/service"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Login(s *service.LoginService) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading login body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loginRequest := &service.LoginRequest{}

		err = json.Unmarshal(body, loginRequest)
		if err != nil {
			log.Printf("Error unmarshalling login body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

		loginResp, err := s.Login(loginRequest, ctx)
		if err != nil {
			log.Printf("Error logging in: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    loginResp.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true, // set to true in production (requires HTTPS)
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int((7 * 24 * time.Hour).Seconds()), // 7 days
		})

		respBody, err := json.Marshal(loginResp)
		if err != nil {
			log.Printf("Error marshalling login response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(respBody)
	}
}
