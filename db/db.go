package db

import (
	"URL-Shortener/model"
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

const col = "Shortener"
const Url = "http://localhost:8080/url/"

type MongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongourl string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancle()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongourl))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	return client, err
}

func NewMongoRepo(mongourl string, mongoDB string, mongoTimeout int) (*MongoRepository, error) {

	repo := &MongoRepository{
		database: mongoDB,
		timeout:  time.Duration(mongoTimeout) * time.Second,
	}

	client, err := newMongoClient(mongourl, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "db.mongorepository")
	}

	repo.client = client
	return repo, nil
}

func (r *MongoRepository) getCollection() (context.Context, *mongo.Collection) {
	ctx, _ := context.WithTimeout(context.Background(), r.timeout)

	collection := r.client.Database(r.database).Collection(col)
	return ctx, collection

}

func (r *MongoRepository) Find(code string) (*model.Redirect, error) {
	ctx, collection := r.getCollection()
	redirect := &model.Redirect{}
	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(errors.New("Not Found"), "db.mongo.find")
		}
		return nil, errors.Wrap(err, "db.mongo.find")
	}

	return redirect, nil
}

func (r *MongoRepository) Store(redirect *model.Redirect) error {

	ctx, collection := r.getCollection()
	filter := bson.M{"url": redirect.URL}

	var temp []model.Redirect
	res, _ := collection.Find(ctx, filter)
	_ = res.All(ctx, &temp)

	if len(temp) > 0 {
		return errors.New("Url already presented with code: " + temp[0].Code)
	}

	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	})

	if err != nil {
		return errors.Wrap(err, "db.mongo.store")
	}

	return nil
}

func (r *MongoRepository) All() ([]model.Redirect, error) {
	// ctx, cancle := context.WithTimeout(context.Background(), r.timeout)
	// defer cancle()
	// collection := r.client.Database(r.database).Collection("Shortener")

	ctx, collection := r.getCollection()
	filter := bson.M{}

	var data []model.Redirect
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("Error Internal Server Error")
	}
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, errors.New("Error Internal Server Error")
	}
	for cursor.Next(context.TODO()) {
		var temp model.Redirect
		err = cursor.Decode(&temp)
		if err != nil {
			return nil, errors.New("Error Internal Server Error")
		}
		temp.Code = Url + temp.Code
		data = append(data, temp)
	}
	if err = cursor.Err(); err != nil {
		return nil, errors.New(err.Error())
	}

	return data, nil
}
