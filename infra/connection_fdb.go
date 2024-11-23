package infra

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type FirestoreHandler struct {
	Client     *firestore.Client
	Collection *firestore.CollectionRef
}

const (
	projectID        = "test_norealtime"
	collectionFbName = "datas"
)

func NewFirestoreHandler() (*FirestoreHandler, error) {
	// Carregar as credenciais
	opt := option.WithCredentialsFile("privt_test_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Firestore")

	collection := client.Collection(collectionFbName)

	return &FirestoreHandler{
		Client:     client,
		Collection: collection,
	}, nil
}

func (handler *FirestoreHandler) Close() {
	if err := handler.Client.Close(); err != nil {
		log.Printf("Error closing Firestore client: %v", err)
	}
}
