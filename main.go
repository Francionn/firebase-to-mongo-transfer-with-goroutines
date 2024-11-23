package main

import (
	"configs/infra"
	"configs/services"
	"log"
	"time"
)

func main() {

	fireHandler, err := infra.NewFirestoreHandler()
	services.CheckError(err, "Error initializing Firestore")
	defer fireHandler.Close()

	mongoHandler, err := infra.NewMongoDBHandler()
	services.CheckError(err, "Error connecting to MongoDB")
	defer mongoHandler.Close()

	firestoreCollection := fireHandler.Collection

	firestoreClient := fireHandler.Client

	mongoCollection := mongoHandler.Collection

	startTime := time.Now()

	err = services.TransferData(firestoreClient, firestoreCollection, mongoCollection)
	services.CheckError(err, "Error transferring data")

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	log.Printf("Data transfer completed in: %v\n", duration)
	log.Println("Data transfer completed successfully")
}
