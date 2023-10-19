package mongo

import (
	"context"
	"log"
	"time"

	"github.com/yellowpuki/tg-bath-bot/internal/er"
	"github.com/yellowpuki/tg-bath-bot/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	records Records
}

type Records struct {
	*mongo.Collection
}

type Record struct {
	EventToken string    `bson:"event-token"`
	UserName   string    `bson:"user-name"`
	CreatedAt  time.Time `bson:"created-at"`
}

// New ...
func New(ctx context.Context, connectString string, connectTimeout time.Duration) Storage {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))
	if err != nil {
		log.Fatalf("can't connect to database: %s", err)
	}

	records := Records{
		Collection: client.Database("storage").Collection("records"),
	}

	return Storage{
		records: records,
	}
}

// Save ...
func (s Storage) Save(ctx context.Context, r *storage.Record) error {
	_, err := s.records.InsertOne(ctx, r)

	if err != nil {
		return er.Wrap("can't save record", err)
	}

	return nil
}

// IsExist ...
func (s Storage) IsExist(ctx context.Context, h string) (bool, error) {
	var result storage.Record

	err := s.records.Collection.FindOne(ctx, bson.M{"eventtoken": h}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, storage.ErrNoRecords
		}
		return false, er.Wrap("can't check is record exists", err)
	}

	return true, nil
}

func (s Storage) LastVisit(ctx context.Context, userName string) (t time.Time, err error) {
	var result storage.Record

	opts := options.FindOne().SetSort(bson.D{{Key: "createdat", Value: 1}})

	err = s.records.Collection.FindOne(ctx, bson.D{{Key: "username", Value: userName}}, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return t, storage.ErrNoRecords
		}
	}
	log.Println(result.CreatedAt)

	return result.CreatedAt, nil
}
