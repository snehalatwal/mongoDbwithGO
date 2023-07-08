package router

import (
	"github.com/gorilla/mux"
	"github.com/snehaMongoDb/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/movies", controller.GetALLMyMovies).Methods("GET")
	router.HandleFunc("/api/movie", controller.InsertMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controller.MarkMovieWatched).Methods("PUT")
	router.HandleFunc("/api/movie/{id}", controller.DeleteAMovie).Methods("DELETE")
	router.HandleFunc("/api/movies", controller.DeleteAllMovies).Methods("DELETE")

	return router
}
