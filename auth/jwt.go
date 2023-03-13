package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"jwt-split-token/database"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// generate a jwt token when user login
func GenerateJWT(password string, username string, accessTokenMaxAge int, jwtKey []byte) (tokenSign string, err error) {

	issueTime := time.Now()
	expirationTime := time.Now().Add(time.Duration(accessTokenMaxAge) * time.Minute)

	// payload
	claims := &JWTClaim{
		Password: password,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issueTime.Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "http://example.com",
		},
	}

	// header + payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// header + payload + signature - convert "header + payload" to string
	headerPayloadString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	// signing with "jwtKey"
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("\ntokenString: ", tokenString) // print just for demonstration!

	// extract "signature"
	tokenSign = strings.Split(tokenString, ".")[2]

	// hashed "signature"
	hashedTokenSign := sha256.Sum256([]byte(tokenSign))

	// store "header + payload" in redis with key "hashedTokenSign"
	redisClient := database.NewRedisClient("redis:6379", "", 0)
	ctx := context.Background()

	// convert hashedTokenSign to string
	hashedTokenSignString := string(hashedTokenSign[:])

	err = redisClient.Set(ctx, hashedTokenSignString, headerPayloadString, 0).Err()
	if err != nil {
		panic(err)
	}

	return tokenSign, err
}

// verify jwt token when users send request to time.ir
func ValidateToken(signedToken string, jwtKey []byte) (jwtToken string, err error) {

	// check if token is not empty
	if signedToken == "" {
		err = errors.New("token is empty")
		return
	}

	// check if token is complete and not just a signature
	if strings.Contains(signedToken, ".") {
		err = errors.New("token is complete")
		return
	}

	// calculate hash of token
	hashedTokenSign := sha256.Sum256([]byte(signedToken))

	// convert hashedTokenSign to string
	hashedTokenSignString := string(hashedTokenSign[:])

	// create new redis client
	redisClient := database.NewRedisClient("redis:6379", "", 0)

	ctx := context.Background()
	// get header + payload from redis
	headerPayloadString, err := redisClient.Get(ctx, hashedTokenSignString).Result()
	if err != nil {
		// checking if token is valid
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
		err = errors.New("token is invalid")
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	// check if token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	// check if token issuer is valid
	if claims.Issuer != "http://example.com" {
		err = errors.New("invalid issuer")
		return
	}

	// return jwt token from "header + payload + signature"
	jwtToken = headerPayloadString + "." + signedToken

	return jwtToken, nil
}
