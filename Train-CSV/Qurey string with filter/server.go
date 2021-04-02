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

type Data struct {
	ID string `bson:"id"`
	// ID        primitive.ObjectID `bson:"id string"`
	TrainNo   string `bson:"trainNumber"`
	TrainName string `bson:"trainName"`
	SEQ       string `bson:"sequence"`
	Code      string `bson:"code"`
	StName    string `bson:"stationName"`
	ATime     string `bson:"arrivalTime"`
	DTime     string `bson:"destinationTime"`
	Distance  string `bson:"distance"`
	SS        string `bson:"sourceStation"`
	SSname    string `bson:"sourceStationName"`
	Ds        string `bson:"destinationStation"`
	DsName    string `bson:"destinationStationName"`
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func getEnv(key string) string {
	err := godotenv.Load(".env")
	CheckError(err)
	value := os.Getenv(key)
	return value
}

func DBConn() *mongo.Client {
	dataBaseUri := getEnv("DATABASE_URI")
	clientOptions := options.Client().ApplyURI(dataBaseUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	CheckError(err)
	err = client.Ping(context.TODO(), nil)
	CheckError(err)
	return client
}
func ReadCsv(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	databaseName := getEnv("DATABASE")
	collectionName := getEnv("COLLECTION")
	collection := db.Database(databaseName).Collection(collectionName)
	defer db.Disconnect(context.TODO())
	fileName := getEnv("FILE_NAME")
	file, err := os.Open(fileName)
	CheckError(err)
	defer file.Close()
	entries, err := csv.NewReader(file).ReadAll()
	CheckError(err)
	var trainData []Data
	id := 0
	for _, lines := range entries {

		line := Data{
			ID:        strconv.Itoa(id),
			TrainNo:   lines[0],
			TrainName: lines[1],
			SEQ:       lines[2],
			Code:      lines[3],
			StName:    lines[4],
			ATime:     lines[5],
			DTime:     lines[6],
			Distance:  lines[7],
			SS:        lines[8],
			SSname:    lines[9],
			Ds:        lines[10],
			DsName:    lines[11],
		}
		_, err := collection.InsertOne(context.TODO(), line)
		CheckError(err)
		trainData = append(trainData, line)
		id++
	}

}
func Fetch(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	databaseName := getEnv("DATABASE")
	collectionName := getEnv("COLLECTION")
	collection := db.Database(databaseName).Collection(collectionName)
	defer db.Disconnect(context.TODO())

	trainNo := r.URL.Query().Get("tNo")
	trainName := r.URL.Query().Get("tName")
	stationName := r.URL.Query().Get("stName")
	var filter bson.D
	if trainNo != "" {
		if len(trainNo) > 0 {
			filter = append(filter, bson.E{"trainNumber", trainNo})

		}
	}

	if trainName != "" {
		if len(trainName) > 0 {
			filter = append(filter, bson.E{"trainName", trainName})

		}
	}

	if stationName != "" {
		if len(stationName) > 0 {
			filter = append(filter, bson.E{"stationName", stationName})

		}
	}
	// filter := bson.M{"trainNumber": trainNo, "trainName": trainName, "stationName": stationName}

	entries, err := collection.Find(context.TODO(), filter)

	CheckError(err)
	var details []Data

	for entries.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var information Data
		err := entries.Decode(&information)
		CheckError(err)

		details = append(details, information)
	}

	pBytes, err := json.Marshal(details)
	CheckError(err)
	w.Write([]byte(pBytes))
}
func short(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	databaseName := getEnv("DATABASE")
	collectionName := getEnv("COLLECTION")
	collection := db.Database(databaseName).Collection(collectionName)
	defer db.Disconnect(context.TODO())
	station1 := r.URL.Query().Get("st1")
	station2 := r.URL.Query().Get("st2")
	filter := bson.M{"stationName": station1, "sourceStationName": station2}
	entries, err := collection.Find(context.TODO(), filter)

	CheckError(err)
	var details []Data

	for entries.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var information Data
		err := entries.Decode(&information)
		// fmt.Println(information)
		CheckError(err)

		details = append(details, information)
	}

	pBytes, err := json.Marshal(details)
	CheckError(err)
	w.Write([]byte(pBytes))

}
func Read(fileName string) {
	db := DBConn()
	databaseName := getEnv("DATABASE")
	collectionName := getEnv("COLLECTION")
	collection := db.Database(databaseName).Collection(collectionName)
	defer db.Disconnect(context.TODO())
	fileName = getEnv("FILE_NAME")
	file, err := os.Open(fileName)
	CheckError(err)
	defer file.Close()
	entries, err := csv.NewReader(file).ReadAll()
	CheckError(err)
	var trainData []Data
	for _, lines := range entries {
		line := Data{
			TrainNo:   lines[0],
			TrainName: lines[1],
			SEQ:       lines[2],
			Code:      lines[3],
			StName:    lines[4],
			ATime:     lines[5],
			DTime:     lines[6],
			Distance:  lines[7],
			SS:        lines[8],
			SSname:    lines[9],
			Ds:        lines[10],
			DsName:    lines[11],
		}
		_, err := collection.InsertOne(context.TODO(), line)
		CheckError(err)
		trainData = append(trainData, line)
	}
}
func main() {
	read := flag.Bool("read", false, "Insert data to database")
	flag.Parse()
	if *read {
		//Getting credentials from Env file
		fileName := getEnv("FILE_NAME")
		Read(fileName)
		fmt.Println("Flag started")
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/read", ReadCsv)
	http.HandleFunc("/fetch", Fetch)
	http.HandleFunc("/short", short)
	fmt.Println("Server started at 9000")
	http.ListenAndServe(":9000", nil)

}

/*-----------Code for BSON MAP-----------
cursor, err := collection.Find(context.TODO(), bson.M{})
	CheckError(err)

	var episodes []bson.M
	if err = cursor.All(context.TODO(), &episodes); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodes)

----------------End-----------------------*/
/*filter := bson.M{
	"$and": []bson.M{
		{"trainNumber": trainNo},
		{"trainName": trainName},
		{"stationName": stationName},
	},
}
*/
//filter := bson.D{{"name", "Ash"}}
//update := bson.D{
//     {"$inc", bson.D{
//         {"age", 1},
//     }},
// }
//updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
