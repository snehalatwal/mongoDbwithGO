package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/snehaMongoDb/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://snehalatwal:WQT2OYwXMQ9ttgpI@cluster0.w85kc4w.mongodb.net/?retryWrites=true&w=majority"
const dbName = "netflix"
const colName = "watchList"

// create collection-very important

var collection *mongo.Collection

// connect to mongo

func init() {

	// client options

	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to mongo

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongoDb")

	// collections

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance created")

}

func insertOneMovie(movie model.Netflix) {

	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The id of movie is: ", inserted.InsertedID)
}

func updateOneMovie(movieId string) {

	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count", result.ModifiedCount)

}

func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Delete count is: ", deleteCount.DeletedCount)
}

func deleteManyMovies() int64 {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	return deleteResult.DeletedCount
}

func getAllMovies() []primitive.M {

	// bson.D is ordered and bson.M unordered
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M
	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	defer cursor.Close(context.Background())
	return movies
}

// func getMovie() []primitive.M {
// 	cur, _ := collection.Find(context.Background(), bson.D{{}})
// 	var movies []primitive.M
// 	for cur.Next(context.Background()) {
// 		var movie bson.M
// 		err := cur.Decode(&movie)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		movies = append(movies, movie)
// 	}

// 	defer cur.Close(context.Background())
// 	return movies
// }

// Actual controller-File

func GetALLMyMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func InsertMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Method", "POST")
	var movie model.Netflix

	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkMovieWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Method", "PUT")
	params := mux.Vars(r)
	updateOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-method", "DELETE")
	params := mux.Vars(r)

	deleteOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "aplication/json")
	w.Header().Set("Allow-Control-Allow-method", "DELETE")

	count := deleteManyMovies()
	json.NewEncoder(w).Encode(count)

}
