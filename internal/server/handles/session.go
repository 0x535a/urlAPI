package handles

import (
	"log"
	"net/http"
	"urlAPI/internal/model"
	"urlAPI/internal/op"
	"urlAPI/util"

	"github.com/gin-gonic/gin"
)

func SessionHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var request op.Session
	if err := c.ShouldBind(&request); err != nil { // auth Error
		util.ErrorPrinter(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authSession := model.Session{
		Token: c.Request.Header.Get("Authorization"),
		Term:  request.LoginTerm,
	}
	request.SessionIP = c.ClientIP()
	request.SessionToken = c.Request.Header.Get("Authorization")

	response, err := op.HandleSession(request, authSession)
	if err != nil {
		log.Printf("%s from %s\n", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
