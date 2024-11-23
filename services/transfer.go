package services

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Transfer data Firestore --> MongoDB
func TransferData(firestoreClient *firestore.Client, firestoreCollection *firestore.CollectionRef, mongoCollection *mongo.Collection) error {
	// Concurrent channel to handle data
	dataChan := make(chan map[string]interface{}, 100)

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Goroutine to fetch data from Firestore
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchDataFromFirestore(firestoreCollection, dataChan); err != nil {
			log.Printf("Error fetching data from Firestore: %v", err)
		}
		close(dataChan)
	}()

	// goroutine to save data to MongoDB in batches
	wg.Add(1)
	go func() {
		defer wg.Done()
		saveDataToMongoInBatches(mongoCollection, dataChan)
	}()

	wg.Wait()
	return nil
}

// fetches data from Firestore and sends it to the channel
func fetchDataFromFirestore(firestoreCollection *firestore.CollectionRef, dataChan chan map[string]interface{}) error {
	ctx := context.Background()

	// get all documents from the Firestore collection
	iter := firestoreCollection.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		item := doc.Data()
		item["id"] = doc.Ref.ID
		dataChan <- item
	}
	return nil
}

// saves data to MongoDB in batches of 500
func saveDataToMongoInBatches(collection *mongo.Collection, dataChan chan map[string]interface{}) {
	var batch []interface{}
	batchSize := 500

	for item := range dataChan {
		batch = append(batch, bson.M(item))

		if len(batch) >= batchSize {
			_, err := collection.InsertMany(context.Background(), batch)
			if err != nil {
				log.Printf("Error inserting batch into MongoDB: %v", err)
			}
			batch = nil
		}
	}

	if len(batch) > 0 {
		_, err := collection.InsertMany(context.Background(), batch)
		if err != nil {
			log.Printf("Error inserting remaining batch into MongoDB: %v", err)
		}
	}
}
