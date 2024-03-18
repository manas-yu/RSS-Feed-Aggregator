package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/manas-yu/rssagg/internal/database"
)

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string
		Url  string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedToFeed(feed))
}

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
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
		feeds, err := cfg.DB.GetFeed(r.Context())
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
			return
		}

		feedJSON, err := json.Marshal(databaseFeedsToFeed(feeds))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to encode feed")
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", feedJSON)

		// Flush the data immediately instead of buffering it for later.
		flusher.Flush()

		// Sleep for a while before the next iteration to avoid high CPU usage.
		time.Sleep(time.Second)
	}
}
