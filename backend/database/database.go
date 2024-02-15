package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"btpTracker/backend/request"
)

var ErrNoDocuments = mongo.ErrNoDocuments

var (
	Client   *mongo.Client
	Database *mongo.Database
)

type DbRow struct {
	ISIN          string    `json:"ISIN" bson:"ISIN"`
	Description   string    `json:"Description" bson:"Description"`
	Last          string    `json:"Last" bson:"Last"`
	Cedola        string    `json:"Cedola" bson:"Cedola"`
	Expiration    string    `json:"Expiration" bson:"Expiration"`
	InsertionDate time.Time `json:"InsertionDate" bson:"InsertionDate"`
}

func Established_connection() (*mongo.Client, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	log.Println("Loaded the .env file")

	uri := os.Getenv("MONGODB_URI")
	username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}
	log.Printf("URI found: %s\n", uri)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: username,
		Password: password,
	}))

	return client, err
}

func CreateDatabase(name string) *mongo.Database {
	log.Println("Enter")

	db := Client.Database(name)
	log.Println(name)
	log.Println(db)
	return db
}

func Create_index(collection *mongo.Collection, index_column string) error {

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: index_column, Value: -1}},
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}

func IsPresent(collectionName string, column string, value string) (*request.Paper, error) {
	collection := Database.Collection(collectionName)
	resp := collection.FindOne(context.TODO(), bson.D{{Key: column, Value: value}})
	if err := resp.Err(); err != nil {
		return nil, err
	}

	paper := &request.Paper{}
	err := resp.Decode(paper)
	if err != nil {
		return nil, err
	}

	return paper, nil
}

type TableRow struct {
	ISIN        string `json:"ISIN" bson:"ISIN"`
	Description string `json:"Description" bson:"Description"`
	Last        string `json:"Last" bson:"Last"`
	Cedola      string `json:"Cedola" bson:"Cedola"`
	Expiration  string `json:"Expiration" bson:"Expiration"`
}

func GetBtpHistory(id string) ([]bson.M, error) {
	// Your MongoDB database and collection names
	// databaseName := "yourDatabase"
	collectionName := "btp"

	collection := Database.Collection(collectionName)
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "ISIN", Value: id}}}}

	// Print the matchStage
	// fmt.Println(matchStage)

	// sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "$natural", Value: 1}}}}
	// addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{
	// 	{Key: "insertionDate", Value: bson.D{
	// 		{Key: "$toDate", Value: bson.D{
	// 			{Key: "$multiply", Value: bson.A{
	// 				bson.D{{Key: "$toLong", Value: "$_id"}},
	// 				1000,
	// 			}},
	// 		}},
	// 	}},
	// }}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "name", Value: "$InsertionDate"},
		{Key: "value", Value: "$Cedola"},
		// Add more fields as needed
	}}}

	// Aggregation pipeline
	pipeline := mongo.Pipeline{matchStage, projectStage}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	// return results
	for _, result := range results {
		fmt.Println(result)
	}
	return results, nil
}
func GetBotHistory(id string) ([]bson.M, error) {
	// Your MongoDB database and collection names
	// databaseName := "yourDatabase"
	collectionName := "bot"

	collection := Database.Collection(collectionName)
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "ISIN", Value: id}}}}

	// Print the matchStage
	// fmt.Println(matchStage)

	// sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "$natural", Value: 1}}}}
	// addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{
	// 	{Key: "insertionDate", Value: bson.D{
	// 		{Key: "$toDate", Value: bson.D{
	// 			{Key: "$multiply", Value: bson.A{
	// 				bson.D{{Key: "$toLong", Value: "$_id"}},
	// 				1000,
	// 			}},
	// 		}},
	// 	}},
	// }}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "name", Value: "$InsertionDate"},
		{Key: "value", Value: "$Last"},
		// Add more fields as needed
	}}}

	// Aggregation pipeline
	pipeline := mongo.Pipeline{matchStage, projectStage}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	// return results
	for _, result := range results {
		fmt.Println(result)
	}
	return results, nil
}
func Insert_element(collectionName string, got any) error {
	collection := Database.Collection(collectionName)
	log.Println((got))
	_, err := collection.InsertOne(context.TODO(), got)
	return err
}

// Returns all the papers that have `paperId` as ancestor.
//
// **Note**: This function is not currently used becuase MongoDB is in the same
// network as the backend. This means that querying the db is quite fast, thus
// we don't need to precompute everything. When the db will be placed in another
// network, it might become useful.
// func GetSubtree(collectionName string, paperId string) (map[string]request.Paper, error) {
// 	collection := Database.Collection(collectionName)

// 	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "Ancestors", Value: paperId}})
// 	if err != nil {
// 		return map[string]request.Paper{}, err
// 	}

// 	defer cursor.Close(context.TODO())

// 	var batchSize int32 = 1024
// 	cursor.SetBatchSize(batchSize)

// 	paper_map := make(map[string]request.Paper)
// 	start := time.Now()
// 	for cursor.Next(context.TODO()) {
// 		var tmp_paper request.Paper
// 		if err := cursor.Decode(&tmp_paper); err != nil {
// 			log.Printf("Error in decode: %s", err)
// 			continue
// 		}
// 		paper_map[tmp_paper.PaperId] = tmp_paper
// 	}
// 	for cursor.RemainingBatchLength() != 0 {
// 		var tmp_paper request.Paper
// 		if err := cursor.Decode(&tmp_paper); err != nil {
// 			log.Printf("Error in decode: %s", err)
// 			continue
// 		}
// 		paper_map[tmp_paper.PaperId] = tmp_paper
// 	}
// 	elapsed := time.Since(start)
// 	log.Printf("Map population: %s\n", elapsed)

// 	return paper_map, nil
// }
