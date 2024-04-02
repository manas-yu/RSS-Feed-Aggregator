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

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	usernameStr := chi.URLParam(r, "username")

	user, err := cfg.DB.GetUserByName(r.Context(), usernameStr)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))

}
func (cfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, dbUser database.User) {
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
		posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
			UserID: dbUser.ID,
			Limit:  10,
		})
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Couldn't get posts")
			return
		}

		postsJSON, err := json.Marshal(databasePostsToPosts(posts))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to encode posts")
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", postsJSON)

		// Flush the data immediately instead of buffering it for later.
		flusher.Flush()

		// Sleep for a while before the next iteration to avoid high CPU usage.
		time.Sleep(time.Second)
	}
}
