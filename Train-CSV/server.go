package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Info struct {
	ID        string `json:"id" bson:"id"`
	TrainNo   string `json:"trainNo" bson:"trainNo"`
	TrainName string `json:"trainName" bson:"trainName"`
	Starts    string `json:"starts" bson:"starts"`
	Ends      string `json:"ends" bson:"ends"`
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func DBConn() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	// defer client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Connected to MongoDB!")
	// collection := client.Database("test").Collection("trainers")
	return client

}

// func connection(client *mongo.Client, database, collection string) *mongo.Collection {

// }
func ReadCsv(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	collection := db.Database("test").Collection("trainers")
	defer db.Disconnect(context.TODO())
	// Open CSV file
	f, err := os.Open("All_Indian_Trains.csv")
	CheckError(err)
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	CheckError(err)
	var dt []Info
	for _, line := range lines {
		data := Info{
			ID:        line[0],
			TrainNo:   line[1],
			TrainName: line[2],
			Starts:    line[3],
			Ends:      line[4],
		}
		_, err = collection.InsertOne(context.TODO(), data)
		CheckError(err)
		dt = append(dt, data)
		fmt.Println(data.ID + " " + data.TrainNo + " " + data.TrainName + " " + data.Starts + " " + data.Ends)
	}

}
func fetchAll(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	collection := db.Database("test").Collection("trainers")
	defer db.Disconnect(context.TODO())
	entries, _ := collection.Find(context.TODO(), bson.D{{}}, options.Find())
	var info []Info
	for entries.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var information Info
		err := entries.Decode(&information)
		CheckError(err)

		info = append(info, information)
	}
	pBytes, err := json.Marshal(info)
	CheckError(err)
	// fmt.Println(pBytes)
	w.Write([]byte(pBytes))
}
func Read(fileName string) {
	db := DBConn()
	collection := db.Database("train").Collection("information")
	defer db.Disconnect(context.TODO())
	file, err := os.Open(fileName)
	lines, err := csv.NewReader(file).ReadAll()
	var dt []Info
	for _, line := range lines {
		data := Info{
			ID:        line[0],
			TrainNo:   line[1],
			TrainName: line[2],
			Starts:    line[3],
			Ends:      line[4],
		}
		_, err = collection.InsertOne(context.TODO(), data)
		CheckError(err)
		dt = append(dt, data)
		fmt.Println(data.ID + " " + data.TrainNo + " " + data.TrainName + " " + data.Starts + " " + data.Ends)
	}

	CheckError(err)

}
func main() {
	// var flagVal bool
	// flag.Var(&flagVal, "name", "help message for flagname")
	// var flagvar int
	// func init() {
	// flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	// }
	read := flag.Bool("read", false, "Insert data to database")
	flag.Parse()
	if *read {
		Read("All_Indian_Trains.csv")
		fmt.Println("Flag started")
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/read", ReadCsv)
	// ReadCsv()
	http.HandleFunc("/fetch", fetchAll)
	fmt.Println("Server started at 8000")
	http.ListenAndServe(":8000", nil)
}
