package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// it is called when a user sends an access request to time.ir
func Access(context *gin.Context) {
	fmt.Println("Accessing time.ir")
	context.Redirect(http.StatusMovedPermanently, "http://www.time.ir/")
}
