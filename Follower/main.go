package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"follower.xws.com/handler"
	"follower.xws.com/repository"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[follower-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[follower-store] ", log.LstdFlags)

	store, err := repository.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseDriverConnection(timeoutContext)
	store.CheckConnection()
	FollowersHandler := handler.NewFollowersHandler(logger, store)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/createFollower", FollowersHandler.CreateFollowing).Methods("POST")

	router.HandleFunc("/getFollowings/{userId}", FollowersHandler.GetFollowingsForUser).Methods("GET")

	router.HandleFunc("/getFollowers/{userId}", FollowersHandler.GetFollowersForUser).Methods("GET")

	router.HandleFunc("/getRecommended/{userId}", FollowersHandler.Recommendations).Methods("GET")

	router.HandleFunc("/removeFollower", FollowersHandler.Unfollow).Methods("DELETE")

	permittedHeaders := handlers.AllowedHeaders([]string{"Requested-With", "Content-Type", "Authorization"})
	permittedOrigins := handlers.AllowedOrigins([]string{"*"})
	permittedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8084", handlers.CORS(permittedHeaders, permittedOrigins, permittedMethods)(router)))
}
