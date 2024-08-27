package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/linus5304/chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	authorId := -1
	authorIdString := r.URL.Query().Get("author_id")
	if authorIdString != "" {
		authorId, err = strconv.Atoi(authorIdString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorId != -1 && dbChirp.AuthorId != authorId {
			continue
		}
		chirps = append(chirps, Chirp{
			Id:       dbChirp.Id,
			Body:     dbChirp.Body,
			AuthorId: dbChirp.AuthorId,
		})
	}

	sortDirection := "asc"
	sortDirectionParam := r.URL.Query().Get("sort")
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].Id > chirps[j].Id
		}
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) handleChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, database.Chirp{
		Id:       dbChirp.Id,
		Body:     dbChirp.Body,
		AuthorId: dbChirp.AuthorId,
	})
}
