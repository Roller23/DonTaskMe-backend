package routing

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
)

type loginData struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type errorMsg struct {
	Message string `json:"message"`
}

func register(c *gin.Context) {
	var ld loginData
	if err := c.ShouldBindJSON(&ld); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//TODO: check if user exists
	//if !isPasswordValid(ld.Password) {
	//	c.JSON(
	//		http.StatusNotAcceptable,
	//		errorMsg{ Message: "Password does not meet minimal rules: eight characters, one digit, one capital letter, one special character, one lowercase letter" },
	//	)
	//}

	_, err := bcrypt.GenerateFromPassword([]byte(ld.Password), 20)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, err.Error())
		return
	}

	//TODO: Save to DB
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
