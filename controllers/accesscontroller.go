package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// if an authorized user sends a request to time.ir
func Access(context *gin.Context) {
	fmt.Println("Accessing time.ir")
	context.Redirect(http.StatusMovedPermanently, "http://www.time.ir/")
}
