package database

import (
	"context"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBase struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func New(Database string, Collection string) (*DataBase, error) {

	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	coll := client.Database(Database).Collection(Collection)
	return &DataBase{
		client: client,
		coll:   coll,
	}, nil
}

func (db DataBase) getUniqueResult() string {
	oldResults, _ := db.coll.Distinct(context.TODO(), "result", bson.D{})
	bInOld := true

	bStr := ""
	for bInOld {

		bInOld = false
		b := make([]byte, 5)
		for i := range b {
			b[i] = BytesSymbols[rand.Intn(len(BytesSymbols))]
		}
		bStr = string(b)
		for i := range oldResults {
			if oldResults[i] == bStr {
				bInOld = true
				break
			}
		}
	}
	return bStr
}

func (db DataBase) FindUrl(url string) (bson.M, error) {

	var result bson.M
	filter := bson.D{{Key: "initial", Value: url}}

	err := db.coll.FindOne(
		context.TODO(), filter,
	).Decode(&result)
	return result, err
}

func (db DataBase) InsertLink(url string) (*mongo.InsertOneResult, error) {
	result, err := db.coll.InsertOne(
		context.TODO(),
		bson.D{
			{Key: "initial", Value: url},
			{Key: "counter", Value: 1},
			{Key: "result", Value: db.getUniqueResult()},
		},
	)

	return result, err
}

func (db DataBase) UpdateCounter(initial string, counter int64) error {
	_, err := db.coll.UpdateOne(
		context.TODO(),
		bson.D{{Key: "initial", Value: initial}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "counter", Value: counter}}}},
	)
	return err
}

func (db DataBase) GetByInitial(initial string) (bson.M, error) {
	res, err := db.FindUrl(initial)
	if err != nil {
		data, err := db.InsertLink(initial)

		if err != nil {
			return nil, err
		}

		filter := bson.D{{Key: "_id", Value: data.InsertedID}}

		err = db.coll.FindOne(
			context.TODO(), filter,
		).Decode(&res)

		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (db DataBase) GetByResult(result string) (bson.M, error) {
	var res bson.M

	filter := bson.D{{Key: "result", Value: result}}
	err := db.coll.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
