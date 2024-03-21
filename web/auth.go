package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

var mySigningKey = []byte(os.Getenv("JWT_SECRETE"))

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func auth(newUser User) (string, error) {
	isMatches := isPassMatches(newUser)
	fmt.Println(isMatches)
	if isMatches {
		token, err := createToken(newUser.Username)
		return token, err
	}
	return "", errors.New("password mismatch")
}

func createToken(username string) (string, error) {
	// Create the Claims
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		Issuer:    "test", //todo pass certificate as issuer
		Subject:   username,
		Audience:  []string{"USER"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		log.SetPrefix("Post failed :::: ")
		log.Fatal(err)
		return "", err
	}
	return ss, nil
}

func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	sub := getClaimsSubject(token.Claims)

	switch {
	case token.Valid:
		fmt.Println("You look nice today")
		return sub, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("That's not even a token")
		return "", err
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// Invalid signature
		fmt.Println("Invalid signature")
		return "", err
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		// Token is either expired or not active yet
		fmt.Println("Token is either expired or not active yet")
		return "", err
	default:
		fmt.Println("Couldn't handle this token:", err)
		return "", err
	}
}

func getBearerToken(c *gin.Context) (string, bool) {
	bearer := c.GetHeader("Authorization")
	if len(bearer) == 0 {
		return "Authorization header is missing", false
	}

	token := bearer[len("Bearer "):]

	return token, true
}

func getClaimsSubject(clm jwt.Claims) string {
	sub, err := clm.GetSubject()

	if err != nil {
		return ""
	}
	return sub
}

func hashPwd(pwd string) (string, error) {
	hsPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hsPwd), nil
}

func isPassMatches(newUser User) bool {
	hashedPwd := users[newUser.Username].Password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(newUser.Password))
	if err != nil {
		return false
	}
	return true
}
