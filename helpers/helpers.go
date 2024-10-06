package helpers

import (
	"fmt"
	"log"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func LoadEnv() {
	// err := godotenv.Load("../../../.env")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
}

func CheckValidEmail(mail string) bool {
	res, err := emailverifier.NewVerifier().Verify(mail)
	if err != nil {
		fmt.Println(err)
	}
	return res.Syntax.Valid
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
