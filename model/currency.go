package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Currency struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Currency string             `bson:"currency" json:"currency"`
	Exchange float64            `bson:"exchange" json:"exchange"`
}

type Converter struct {
	Currency string  `bson:"currency,omitempty" json:"currency,omitempty"`
	Amount   float64 `bson:"amount, omitempty" json:"amount,omitempty"`
}
