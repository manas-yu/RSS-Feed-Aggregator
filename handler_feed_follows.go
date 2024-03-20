package main

import (
	"database/sql"
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
		Name:      sql.NullString{String: user.Name, Valid: true},
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) handlerGetUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	// Make sure that the writer supports flushing.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		select {
		case <-r.Context().Done():
			// The client has closed the connection.
			fmt.Println("connection closed in feed follows")
			return
		default:
			// The client is still connected.
		}

		userFeedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Couldn't get feed follows")
			return
		}

		feedFollowsJSON, err := json.Marshal(databaseUserFeedFollowsToFeedFollows(userFeedFollows))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to encode feed follows")
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", feedFollowsJSON)

		// Flush the data immediately instead of buffering it for later.
		flusher.Flush()

		// Sleep for a while before the next iteration to avoid high CPU usage.
		time.Sleep(time.Second)
	}
}

func (cfg *apiConfig) handlerDeleteFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowId")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{ msg string }{msg: "deleted feed follow"})

}
