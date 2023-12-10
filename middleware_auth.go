package main

import (
	"fmt"
	"net/http"

	"github.com/manas-yu/rssagg/internal/auth"
	"github.com/manas-yu/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 400, "could not get api key")
			fmt.Println(err)
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		fmt.Println(apiKey)
		if err != nil {
			respondWithError(w, 400, "could not get user")
			fmt.Println(err)
			return
		}
		handler(w, r, user)
	}
}
