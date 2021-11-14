package routing

import (
	db2 "DonTaskMe-backend/internal/db"
	models2 "DonTaskMe-backend/internal/models"
	"DonTaskMe-backend/pkg/hash"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"regexp"
)

func register(c *gin.Context) {
	var user models2.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if userAlreadyExists(&user.Username) {
		c.JSON(http.StatusConflict, "user with given username already exists")
		return
	}

	//if !isPasswordValid(ld.Password) {
	//	c.JSON(
	//		http.StatusNotAcceptable,
	//		errorMsg{ Message: "Password does not meet minimal rules: eight characters, one digit, one capital letter, one special character, one lowercase letter" },
	//	)
	//}

	hashedPassword, err := hash.Generate(&user.Password)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, err.Error())
		return
	}
	user.Password = hashedPassword

	usersCollection := db2.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("USERS_COLLECTION"))
	_, err = usersCollection.InsertOne(context.TODO(), user)
	c.Status(http.StatusCreated)
}

//TODO: Fix the regexp
func isPasswordValid(password string) bool {
	exp, err := regexp.Compile(`.{8,}(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*\W)`)
	if err != nil {
		log.Fatalln("Regexp did not compile: ", err.Error())
	}
	return exp.Match([]byte(password))
}

func userAlreadyExists(username *string) bool {
	var res models2.User
	usersCollection := db2.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("USERS_COLLECTION"))
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&res)
	return err != mongo.ErrNoDocuments
}
