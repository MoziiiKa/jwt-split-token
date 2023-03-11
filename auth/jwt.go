package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"jwt-split-token/database"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("p3s6v9y$B&E(H+MbQeThWmZq4t7w!z%C")

type JWTClaim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// it is called when user login
func GenerateJWT(password string, username string) (tokenSign string, err error) {

	issueTime := time.Now()
	expirationTime := time.Now().Add(30 * time.Minute)

	// payload
	claims := &JWTClaim{
		Password: password,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issueTime.Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "test",
		},
	}
	fmt.Println("\nclaims: ", claims)

	// header + payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("\ntoken: ", token)

	// header + payload + signature
	// convert "header + payload" to string
	headerPayloadString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("\ntokenString: ", tokenString)

	// extract "signature"
	tokenSign = strings.Split(tokenString, ".")[2]
	fmt.Println("\ntokenSign: ", tokenSign)

	// hashed "signature"
	hashedTokenSign := sha256.Sum256([]byte(tokenSign))
	fmt.Println("\nhashedTokenSign: ", hashedTokenSign)

	// store "header + payload" in a cache with key "hashedTokenSign"
	// create new redis client
	redisClient := database.NewRedisClient("localhost:6379", "", 0)
	ctx := context.Background()
	// convert hashedTokenSign to string
	hashedTokenSignString := string(hashedTokenSign[:])

	err = redisClient.Set(ctx, hashedTokenSignString, headerPayloadString, 0).Err()
	if err != nil {
		panic(err)
	}
	// send "signature" (to user)

	return tokenSign, err
}

// it is called when user sends an access request to time.ir
func ValidateToken(signedToken string) (jwtToken string, err error) {

	// check if token is not empty
	if signedToken == "" {
		err = errors.New("token is empty")
		return
	}

	// check if token is complete and not just a signature
	// if it has dots, it means that it is complete
	if strings.Contains(signedToken, ".") {
		err = errors.New("token is complete")
		return
	}

	// check if token is valid & calculate hash of token
	hashedTokenSign := sha256.Sum256([]byte(signedToken))

	// convert hashedTokenSign to string
	hashedTokenSignString := string(hashedTokenSign[:])

	// create new redis client
	redisClient := database.NewRedisClient("localhost:6379", "", 0)

	ctx := context.Background()
	// get header + payload from redis
	headerPayloadString, err := redisClient.Get(ctx, hashedTokenSignString).Result()
	if err != nil {
		// if there is no such key in redis, it means that token is invalid
		err = errors.New("token is invalid")
		return
	}

	// parse header + payload
	token, err := jwt.ParseWithClaims(
		headerPayloadString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		// if there is an error, it means that token is invalid
		err = errors.New("token is invalid")
		return
	}

	// check if token is expired
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	// return jwt token from header + payload + signature
	jwtToken = headerPayloadString + "." + signedToken

	return jwtToken, nil
}
