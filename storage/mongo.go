package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/danielkraic/knihomol/books"
	"github.com/danielkraic/knihomol/configuration"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Storage application DB
type Storage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

//NewStorage creates new storage
func NewStorage(cfg *configuration.Storage, timeout time.Duration) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo %s: %s", cfg.URI, err)
	}

	return &Storage{
		client:     client,
		collection: client.Database(cfg.DBName).Collection(cfg.CollectionName),
	}, nil
}

//GetBooks retrieve books from DB
func (s *Storage) GetBooks(ctx context.Context) ([]*books.BookDetails, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %s", err)
	}
	defer func() {
		err = cur.Close(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	var result []*books.BookDetails
	err = cur.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode books from DB: %s", err)
	}

	return result, nil
}

//SaveBook saves book from DB
func (s *Storage) SaveBook(ctx context.Context, book *books.BookDetails) error {
	data, err := bson.Marshal(book)
	if err != nil {
		return fmt.Errorf("failed to marshall book to bson: %s", err)
	}

	_, err = s.collection.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to save book: %s", err)
	}

	return nil
}
