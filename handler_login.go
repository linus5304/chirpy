package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/linus5304/chirpy/internal/auth"
)

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params")
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	defaultExpiration := 60 * 60
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	token, err := auth.MakeJWT(user.Id, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token")
	}

	err = cfg.DB.SaveRefreshToken(user.Id, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:    user.Id,
			Email: user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
