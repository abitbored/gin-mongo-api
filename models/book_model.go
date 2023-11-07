package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID       primitive.ObjectID `json:"id,omitempty"`
	Title    string             `json:"title,omitempty" validate:"required"`
	Author   string             `json:"author,omitempty" validate:"required"`
	Synopsis string             `json:"synopsis,omitempty" validate:"required"`
}
