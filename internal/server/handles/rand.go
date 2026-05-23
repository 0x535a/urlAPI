package handles

import (
	"time"
	"urlAPI/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"urlAPI/internal/op"
	"urlAPI/internal/server/middleware"
	"urlAPI/util"
)

func RandHandler(c *gin.Context) {
	target := c.Query("user") + "/" + c.Query("repo")
	referer := c.Request.Referer()
	ip := c.ClientIP()
	task := model.Task{
		UUID:     uuid.New().String(),
		Time:     time.Now(),
		IP:       ip,
		Type:     util.TypeMap["rand"],
		Target:   target,
		Region:   util.GetRegion(ip),
		Referer:  referer,
		Device:   util.GetDeviceType(c.GetHeader("User-Agent")),
		API:      c.Query("api"),
		MoreInfo: c.Query("more"),
	}
	_, result, err := op.GenerateRandom(task, middleware.GetSkipDB(c))
	util.ErrorPrinter(err)
	if err != nil {
		errorReturner(c, err)
		return
	}
	returner(c, result)
}
