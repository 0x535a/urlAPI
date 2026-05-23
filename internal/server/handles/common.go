package handles

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"urlAPI/internal/op"
)

func returner(c *gin.Context, result op.GenerateResult) {
	if c.Query("format") == "json" {
		c.JSON(http.StatusOK, result)
	} else {
		c.Redirect(http.StatusFound, result.URL)
	}
}

func errorReturner(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return `https://`
	}
	if scheme := c.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme + `://`
	}
	return `http://`
}
