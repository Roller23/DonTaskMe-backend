package routing

import (
	"DonTaskMe-backend/internal/db"
	"DonTaskMe-backend/internal/models"
	"DonTaskMe-backend/pkg/hash"
	"context"
	"github.com/gin-gonic/gin"
	nano "github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"regexp"
)

func register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if userAlreadyExists(&user.Username) {
		c.JSON(http.StatusConflict, "user with given username already exists")
		return
	}

	//TODO: validate password
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

	uid, err := nano.Nanoid()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	user.Uid = &uid

	usersCollection := db.Client.Database(db.Name).Collection(db.UsersCollectionName)
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
	var res models.User
	usersCollection := db.Client.Database(db.Name).Collection(db.UsersCollectionName)
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&res)
	return err != mongo.ErrNoDocuments
}
