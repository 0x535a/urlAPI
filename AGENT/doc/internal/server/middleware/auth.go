package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	_error "zhongxin/internal/error"
	"zhongxin/internal/op"
	"zhongxin/util"
)

func GetAuthMiddleware(typs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthMiddleware(c, typs...)
	}
}

func AuthMiddleware(c *gin.Context, typs ...string) {
	token := c.Request.Header.Get("Authorization")
	token, ok := strings.CutPrefix(token, "Bearer ")
	if !ok || len(token) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No Token",
		})
		c.Abort()
		return
	}

	tokenDB, _, errCode := op.GetToken(token)
	switch {
	case errCode == _error.DBRecordNotFound:
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unknown Token",
		})
		c.Abort()
		return
	case util.TimeNow().Unix() > tokenDB.ExpireTime:
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Expired Token",
		})
		c.Abort()
		return
	}

	c.Set("user", tokenDB.UserID)
	c.Set("userType", tokenDB.UserType)
	for _, typ := range typs {
		if typ == tokenDB.UserType {
			c.Next()
			return
		}
	}

	c.JSON(http.StatusForbidden, gin.H{
		"error": "User Permission Denied",
	})
	c.Abort()
}
