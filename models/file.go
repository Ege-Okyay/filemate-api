package models

import (
	"time"
)

type File struct {
	Name      string    `json:"name" bson:"name"`
	Path      string    `json:"path" bson:"path"`
	UserID    string    `json:"userId" bson:"userId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
