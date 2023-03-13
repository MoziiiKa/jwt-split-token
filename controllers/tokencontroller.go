package controllers

import (
	"jwt-split-token/auth"
	"jwt-split-token/database"
	"jwt-split-token/models"
	"os"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// it is called when login
func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// check if email exists
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	// check if password is correct
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	accessTokenMaxAgeStr := os.Getenv("ACCESS_TOKEN_MAX_AGE")

	// convert to int
	accessTokenMaxAge, err := strconv.Atoi(accessTokenMaxAgeStr)
	if err != nil {
		panic(err)
	}

	jwtKeyStr := os.Getenv("JWT_KEY")
	jwtKey := []byte(jwtKeyStr)

	// print only for debuging
	// fmt.Println("jwtKey: ", jwtKey)

	// getting JWT token from jwt.go
	tokenSign, err := auth.GenerateJWT(user.Password, user.Username, accessTokenMaxAge, jwtKey)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send token to user
	context.JSON(http.StatusOK, gin.H{"token": tokenSign})
}
