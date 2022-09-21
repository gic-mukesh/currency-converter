package service

import (
	"context"
	"currency-converter/model"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Currency struct {
	Server     string
	Database   string
	Collection string
}

var Collection *mongo.Collection
var ctx = context.TODO()
var insertDocs int

func (e *Currency) Connect() {
	clientOptions := options.Client().ApplyURI(e.Server)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database(e.Database).Collection(e.Collection)
}

func (e *Currency) Insert(curr model.Currency) error {

	_, err := Collection.InsertOne(ctx, curr)

	if err != nil {
		return errors.New("Unable To Insert New Record")
	}
	return nil
}

func (e *Currency) Convert(converter model.Converter) (float64, error) {
	var currency []*model.Currency

	cur, err := Collection.Find(ctx, bson.D{primitive.E{Key: "currency", Value: converter.Currency}})

	if err != nil {
		return 0, errors.New("Unable TO Apply Query")
	}

	for cur.Next(ctx) {
		var e model.Currency
		err := cur.Decode(&e)
		if err != nil {
			return 0, err
		}
		currency = append(currency, &e)
	}
	if currency == nil {
		return 0, errors.New("No data present in db for given currency")
	}
	convrt := converter.Amount * currency[0].Exchange

	if convrt == 0 {
		return 0, errors.New("Unable TO Convert Given Amount")
	}

	return convrt, nil
}
