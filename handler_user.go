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

	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: dbUser.ID,
		Limit:  10,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts")
		return
	}
	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
