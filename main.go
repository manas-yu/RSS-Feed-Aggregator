package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/manas-yu/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("portString not found")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbURL not found")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("could not connect to db", err)
	}
	queries := database.New(conn)
	apiCfg := apiConfig{DB: queries}
	go startScraping(queries, 10, time.Minute)
	router := chi.NewRouter()
	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300},
	),
	)
	v1Router := chi.NewRouter()
	router.Mount("/v1", v1Router)
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.Get("/users/{username}", apiCfg.handlerGetUser)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedsCreate))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetUserFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFollow))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	srv := &http.Server{Handler: router, Addr: ":" + portString}
	log.Printf("server starting in port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

// val request = Request.Builder()
//     .url("http://your-server.com/feeds")
//     .build()

// val client = OkHttpClient()

// client.newCall(request).execute().use { response ->
//     if (!response.isSuccessful) throw IOException("Unexpected code $response")

//     for (line in response.body!!.source().buffer().readUtf8().split("\n")) {
//         if (line.startsWith("data: ")) {
//             val json = line.substring(6)
//             val feed = parseJsonToFeed(json)  // Implement this function to parse the JSON to your Feed object
//             println(feed)
//         }
//     }
// }
//--------------------------------------------------------------
// val request = Request.Builder()
//     .url("http://your-server.com/v1/posts")
//     .addHeader("Authorization", "Bearer your-token")
//     .build()

// val client = OkHttpClient()

// client.newCall(request).execute().use { response ->
//     if (!response.isSuccessful) throw IOException("Unexpected code $response")

//     for (line in response.body!!.source().buffer().readUtf8().split("\n")) {
//         if (line.startsWith("data: ")) {
//             val json = line.substring(6)
//             val post = parseJsonToPost(json)  // Implement this function to parse the JSON to your Post object
//             println(post)
//         }
//     }
// }
