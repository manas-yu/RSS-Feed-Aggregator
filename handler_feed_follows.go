package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/manas-yu/rssagg/internal/database"
)

func (cfg *apiConfig) handlerFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedId,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) handlerGetUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	userFeedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feed follows")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserFeedFollowsToFeedFollows(userFeedFollows))
}

func (cfg *apiConfig) handlerDeleteFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdString := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdString)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse id")
		return
	}
	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{UserID: user.ID, ID: feedFollowId})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete follow")
		return
	}
	fmt.Println("Deleted follow")
	respondWithJSON(w, 200, struct{}{})

}
