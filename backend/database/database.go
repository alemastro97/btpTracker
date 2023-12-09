package database

import (
	"context"
	"log"
	"os"

	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"citation-graph/backend/request"
)

var ErrNoDocuments = mongo.ErrNoDocuments

var (
	Client   *mongo.Client
	Database *mongo.Database
)

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
	return Client.Database(name)
}

func Create_index(collection *mongo.Collection, index_column string) error {
	// collection := client.Database("connected-papers").Collection("papers")

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

func Insert_element(collectionName string, got any) error {
	collection := Database.Collection(collectionName)

	_, err := collection.InsertOne(context.TODO(), got)
	return err
}

// Returns all the papers that have `paperId` as ancestor.
//
// **Note**: This function is not currently used becuase MongoDB is in the same
// network as the backend. This means that querying the db is quite fast, thus
// we don't need to precompute everything. When the db will be placed in another
// network, it might become useful.
func GetSubtree(collectionName string, paperId string) (map[string]request.Paper, error) {
	collection := Database.Collection(collectionName)

	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "Ancestors", Value: paperId}})
	if err != nil {
		return map[string]request.Paper{}, err
	}

	defer cursor.Close(context.TODO())

	var batchSize int32 = 1024
	cursor.SetBatchSize(batchSize)

	paper_map := make(map[string]request.Paper)
	start := time.Now()
	for cursor.Next(context.TODO()) {
		var tmp_paper request.Paper
		if err := cursor.Decode(&tmp_paper); err != nil {
			log.Printf("Error in decode: %s", err)
			continue
		}
		paper_map[tmp_paper.PaperId] = tmp_paper
	}
	for cursor.RemainingBatchLength() != 0 {
		var tmp_paper request.Paper
		if err := cursor.Decode(&tmp_paper); err != nil {
			log.Printf("Error in decode: %s", err)
			continue
		}
		paper_map[tmp_paper.PaperId] = tmp_paper
	}
	elapsed := time.Since(start)
	log.Printf("Map population: %s\n", elapsed)

	return paper_map, nil
}
