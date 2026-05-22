package handles

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/internal/model"
	"urlAPI/internal/op"
	"urlAPI/util"
)

func SessionHandler(c *gin.Context) {
	session, dbSession, err := sessionBuilder(c)
	if err != nil {
		util.ErrorPrinter(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := session.Process(&dbSession); err != nil {
		log.Printf("%s from %s\n", err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, session)
	}
}

func sessionBuilder(c *gin.Context) (op.Session, model.Session, error) {
	c.Header("Access-Control-Allow-Origin", "*")
	var session op.Session
	if err := c.ShouldBind(&session); err != nil { // auth Error
		return session, model.Session{}, err
	}
	dbSession := model.Session{
		Token: c.Request.Header.Get("Authorization"),
		Term:  session.LoginTerm,
	}
	session.SessionIP = c.ClientIP()
	session.SessionToken = c.Request.Header.Get("Authorization")
	return session, dbSession, nil
}
