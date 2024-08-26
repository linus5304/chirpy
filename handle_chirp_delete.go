package main

import (
	"net/http"
	"strconv"

	"github.com/linus5304/chirpy/internal/auth"
)

func (cfg *apiConfig) handleChirpDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token")
		return
	}

	secret, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	userId, err := strconv.Atoi(secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user id")
		return
	}

	chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't parse chirp Id")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp")
		return
	}

	if chirp.AuthorId != userId {
		respondWithError(w, http.StatusForbidden, "Not authorized to delete this chirp")
		return
	}

	err = cfg.DB.DeleteChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
