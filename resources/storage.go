package resources

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return nil, fmt.Errorf("connect to mongo %s: %s", cfg.URI, err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ping mongo %s: %s", cfg.URI, err)
	}

	return &Storage{
		client:     client,
		collection: client.Database(cfg.DBName).Collection(cfg.CollectionName),
	}, nil
}

//GetBooks retrieve books from DB
func (s *Storage) GetBooks(ctx context.Context) ([]*models.Book, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("get books: %s", err)
	}
	defer func() {
		err = cur.Close(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	result := []*models.Book{}

	for cur.Next(ctx) {
		var item *models.Book
		err = cur.Decode(&item)
		if err != nil {
			return nil, fmt.Errorf("decode book from DB: %s", err)
		}

		result = append(result, item)
	}

	return result, nil
}

//SaveBook saves book to DB
func (s *Storage) SaveBook(ctx context.Context, book *models.Book) error {
	filter := bson.D{primitive.E{Key: "_id", Value: book.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: book}}
	opt := &options.UpdateOptions{}
	opt.SetUpsert(true)

	_, err := s.collection.UpdateOne(ctx, filter, update, opt)
	if err != nil {
		return fmt.Errorf("save book: %s", err)
	}

	return nil
}

//RemoveBook removes book from DB
func (s *Storage) RemoveBook(ctx context.Context, bookID string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: bookID}}

	_, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("remove book: %s", err)
	}

	return nil
}

//RemoveBooks removes books from DB
func (s *Storage) RemoveBooks(ctx context.Context, bookIDs []string) error {
	filter := bson.M{"_id": bson.M{"$in": bookIDs}}

	_, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("remove books: %s", err)
	}

	return nil
}
