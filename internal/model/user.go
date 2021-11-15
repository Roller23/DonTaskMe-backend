package model

import (
	"DonTaskMe-backend/internal/service"
	"context"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Uid      *string `json:"uid,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Token    *string `json:"token,omitempty"`
}

func (u *User) AssignNewToken() error {
	uid, _ := nano.Nanoid()
	u.Token = &uid
	usersCollection := service.DB.Collection(service.UsersCollectionName)
	_, err := usersCollection.UpdateOne(context.TODO(), bson.M{"uid": u.Uid}, bson.D{{"$set", bson.D{{"token", u.Token}}}})
	return err
}
