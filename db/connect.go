package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MDB : Exported MongoDB Instance
var MDB MongoDB = MongoDB{}

// MongoDB ...
type MongoDB struct {
	client     *mongo.Client
	URL        string
	db         *mongo.Database
	collection *mongo.Collection
}

// Isconnected : Check If Connected to Database
func (m MongoDB) Isconnected() bool {
	if m.db == nil {
		return false
	} else {
		return true
	}
}

// PingTest : Test Successful Connection
func (m MongoDB) PingTest() error {
	if m.client == nil {
		er := m.Connect()
		if er != nil {
			return er
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetCollInstance : Get Collection Instance
func (m *MongoDB) GetCollInstance() *mongo.Collection {
	return m.collection
}

//Makes a new connection and returns database instance
//this is explicit and everything must be handled

//TempConnection : Create a temp connection to different database
func (m *MongoDB) TempConnection(db string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if m.URL == "" {
		m.URL = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.URL))
	if err != nil {
		panic(err)
	}

	dbc := client.Database(db)
	return dbc
}

// Connect : Connect to database
func (m *MongoDB) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if m.URL == "" {
		m.URL = "mongodb://localhost:27017"
	}
	var err error
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.URL))
	if err != nil {
		return err
	}

	//add ping here
	err = m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetDatabase : Connect to Database
func (m *MongoDB) GetDatabase(name string) {
	m.db = m.client.Database(name)
}

// GetCollection : Get COllection
func (m *MongoDB) GetCollection(name string) {
	m.collection = m.db.Collection(name)
}

// UpdateDocument : Update Existing Document or create New One
func (m *MongoDB) UpdateDocument(filter interface{}, dat interface{}) (*mongo.UpdateResult, error) {
	// fmt.Println("Update Document called")
	opts := options.Update().SetUpsert(true)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := m.collection.UpdateOne(ctx, filter, dat, opts)
	if err != nil {
		return result, err
	}
	// fmt.Printf("Matchedcount %v , Modified COunt %v, with upsert id %v\n ", result.MatchedCount, result.ModifiedCount, result.UpsertedID)
	return result, nil
}

// InsertOne : Insert One Document
func (m *MongoDB) InsertOne(dat interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := m.collection.InsertOne(ctx, dat)
	if err != nil {
		fmt.Println("GOt error while inserting document")
	} else {
		fmt.Println("Successfully Done")
		fmt.Println(result)
	}

	return err
}

// FindOne : Find One Document using Filter
func (m *MongoDB) FindOne(filter interface{}, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.collection.FindOne(ctx, filter).Decode(data); err != nil {
		return data, err
	}

	return data, nil

}

// FindAll : Find All Possible Matches
func (m *MongoDB) FindAll() ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	var arr []bson.M
	if err = cursor.All(ctx, &arr); err != nil {
		panic(err)
	}

	fmt.Println(arr)

	var x []interface{}

	return x, nil
}

//This will create required collection for our note taker

//InitializeProject : Initalize Project According to Sandman Needs
func (m *MongoDB) InitializeProject() error {
	//create global collection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.db.CreateCollection(ctx, "global")
	if err != nil {
		return err
	}

	//create notes collection
	err = m.db.CreateCollection(ctx, "notes")
	if err != nil {
		return err
	}

	//create checklists collection
	err = m.db.CreateCollection(ctx, "checklists")
	if err != nil {
		return err
	}

	//create tooling  collection
	//tooling collection is used to save tool outputs
	err = m.db.CreateCollection(ctx, "tooling")
	if err != nil {
		return err
	}

	return nil
}

// ValidateProject : Just Check if Files are as intended by application
func (m *MongoDB) ValidateProject() bool {
	dat, err := m.ListDBCollections()
	if err != nil {
		panic(err)
	}
	reqcount := 0
	for _, v := range dat {
		if v == "global" || v == "notes" || v == "tooling" || v == "checklists" {
			reqcount = reqcount + 1
		}
	}
	if reqcount == 4 {
		return true
	} else {
		fmt.Println("Looks Like You have chosen the Wrong Project or Version of it ")
		return false
	}

}

// ListDBCollections : List All Collections of Current Database
func (m *MongoDB) ListDBCollections() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := m.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

// ListDatabases : List all Databases
func (m *MongoDB) ListDatabases() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := m.client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

// DropDatabase : Drop a Database
func (m *MongoDB) DropDatabase(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	m.GetDatabase(name)
	err := m.db.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Disconnect : throws error when fails
func (m *MongoDB) Disconnect() {
	er := m.PingTest()
	if er != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := m.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
