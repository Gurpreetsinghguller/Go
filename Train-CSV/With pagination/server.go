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
	"strconv"

	"github.com/joho/godotenv"
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

var (
	ch = make(chan bool, 10)
	// database_uri = "mongodb://localhost:27017"
	// database     = "test"
	// collection   = "trainers"
	// fileName     = "All_Indian_Trains.csv"
)

func goDotEnv(key string) string {
	err := godotenv.Load(".env")
	CheckError(err)
	return os.Getenv(key)

}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func DBConn() *mongo.Client {

	//Getting credentials from Env file
	dataBaseUri := goDotEnv("DATABASE_URI")

	// Set client options

	clientOptions := options.Client().ApplyURI(dataBaseUri)

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
	//Getting credentials from Env file
	database := goDotEnv("DATABASE_NAME")
	colle := goDotEnv("COLLECTION")
	collection := db.Database(database).Collection(colle)
	defer db.Disconnect(context.TODO())

	//Getting credentials from Env file
	fileName := goDotEnv("FILE_NAME")
	// Open CSV file
	f, err := os.Open(fileName)
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
	//Getting credentials from Env file
	database := goDotEnv("DATABASE_NAME")
	//Getting credentials from Env file
	colle := goDotEnv("COLLECTION")
	collection := db.Database(database).Collection(colle)
	defer db.Disconnect(context.TODO())
	var info []Info

	//Getting query String

	query := r.URL.Query()
	page := query.Get("page")

	if page == "" {
		page = "1"
	}

	const limit = 6
	intpg, _ := strconv.Atoi(page)
	skip := (intpg - 1) * 5
	// var inn int64 = skip
	// if skip == 0 {
	// 	skip = limit
	// }
	fmt.Println(page, skip)
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(int64(skip))

	// filter := bson.M{"limit": bson.M{"skip": 42}}
	entries, _ := collection.Find(context.TODO(), bson.D{{"ENDS", "VISHAKHAPATNAM"}}, findOptions)
	for entries.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var information Info
		err := entries.Decode(&information)
		CheckError(err)

		info = append(info, information)
	}
	entries.Close(context.TODO())
	pBytes, err := json.Marshal(info)
	CheckError(err)
	// fmt.Println(pBytes)
	w.Write([]byte(pBytes))
}
func Read(fileName string) {
	db := DBConn()
	//Getting credentials from Env file
	colle := goDotEnv("FILE_NAME")
	database := goDotEnv("DATABASE_NAME")
	collection := db.Database(database).Collection(colle)
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

	read := flag.Bool("read", false, "Insert data to database")
	flag.Parse()
	if *read {
		//Getting credentials from Env file
		fileName := goDotEnv("FILE_NAME")
		Read(fileName)
		fmt.Println("Flag started")
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/read", ReadCsv)
	// ReadCsv()
	//getProductsHandler := http.HandlerFunc(fetchAll)
	//http.HandleFunc("/fetch", getProductsHandler)
	http.HandleFunc("/fetch", fetchAll)
	fmt.Println("Server started at 9000")
	http.ListenAndServe(":9000", nil)
}

/*
func page(request response, next){
	//Getting query string
	query := r.URL.query()
	// fmt.Println(query)
	page := query.Get("page")

	fmt.Println("page is ......", page)
	if !page {
		page = 1
	}

	const limit = 5
	const skip = (page - 1) * size



}
*/
/*var flagVal bool
flag.Var(&flagVal, "name", "help message for flagname")
var flagvar int
func init() {
flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
}

*/

/*
    //Skipping data from 1 and 2 page On page 3

 const users= await Users.find({} {} {limit: limit, skip:skip})
 res.send(users)
*/
