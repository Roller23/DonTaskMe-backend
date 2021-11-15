package helpers

import (
	"DonTaskMe-backend/internal/model"
	"DonTaskMe-backend/internal/service"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(username *string) (*model.User, error) {
	var res model.User
	usersCollection := service.DB.Collection(service.UsersCollectionName)
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func FindUserByToken(token string) (*model.User, error) {
	var res model.User
	usersCollection := service.DB.Collection(service.UsersCollectionName)
	err := usersCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func UserAlreadyExists(username *string) bool {
	var res model.User
	usersCollection := service.Client.Database(service.Name).Collection(service.UsersCollectionName)
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&res)
	return err != mongo.ErrNoDocuments
}

//TODO: Fix the regexp
//func IsPasswordValid(password string) bool {
//	exp, err := regexp.Compile(`.{8,}(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*\W)`)
//	if err != nil {
//		log.Fatalln("Regexp did not compile: ", err.Error())
//	}
//	return exp.Match([]byte(password))
//}
