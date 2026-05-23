package server

import (
	"runtime/debug"
	"zhongxin/cmd/flags"
	_error "zhongxin/internal/error"
	"zhongxin/internal/server/handles"
	"zhongxin/internal/server/middleware"
	"zhongxin/util"

	"github.com/gin-gonic/gin"
)

func RouterRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				util.Log.Error("Recover panic: ", err.(error).Error())
				debug.PrintStack()
				c.JSON(500, gin.H{
					"code": _error.ServerOtherError,
					"msg":  err.(error).Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func Init(e *gin.Engine) {
	e.Use(RouterRecovery())
	v1 := e.Group("/v1")

	user := v1.Group("/user")
	user.POST("/login", handles.UserLoginHandler)
	user.POST("/bind", handles.UserBindHandler)
	if flags.Dev {
		user.GET("/unbind", handles.UserUnbindHandler)
	}

	workPeriod := v1.Group("/workPeriod") // done except attention and notification
	workPeriod.POST("/start", middleware.GetAuthMiddleware("worker"), handles.WorkPeriodStartHandler)
	workPeriod.POST("/stop", middleware.GetAuthMiddleware("worker"), handles.WorkPeriodStopHandler)
	workPeriod.GET("/get", middleware.GetAuthMiddleware("admin", "worker"), handles.WorkPeriodGetHandler)
	workPeriod.PATCH("/update", middleware.GetAuthMiddleware("admin"), handles.WorkPeriodUpdateHandler)
	workPeriod.DELETE("/delete", middleware.GetAuthMiddleware("admin"), handles.WorkPeriodDeleteHandler)
	workPeriod.GET("/list", middleware.GetAuthMiddleware("admin", "worker"), handles.WorkPeriodListHandler)

	worker := v1.Group("/worker") // done except salary details
	worker.GET("/getSalaryDetails", middleware.GetAuthMiddleware("admin", "worker"), handles.WorkerSalaryDetailHandler)
	worker.GET("/getSalaryOverview", middleware.GetAuthMiddleware("admin"), handles.WorkerSalaryOverviewHandler)
	worker.GET("/getSalaryExcel", middleware.GetAuthMiddleware("admin"), handles.WorkerSalaryExcelHandler)
	worker.GET("/get", middleware.GetAuthMiddleware("admin", "worker"), handles.WorkerGetHandler)
	worker.GET("/list", middleware.GetAuthMiddleware("admin"), handles.WorkerListHandler)
	worker.DELETE("/delete", middleware.GetAuthMiddleware("admin"), handles.WorkerDeleteHandler)
	worker.POST("/add", middleware.GetAuthMiddleware("admin"), handles.WorkerAddHandler)
	worker.PATCH("/update", middleware.GetAuthMiddleware("admin"), handles.WorkerUpdateHandler)

	machine := v1.Group("/machine") // done
	machine.GET("/get", middleware.GetAuthMiddleware("admin", "worker"), handles.MachineGetHandler)
	machine.GET("/list", middleware.GetAuthMiddleware("admin", "worker"), handles.MachineListHandler)
	machine.PATCH("/update", middleware.GetAuthMiddleware("admin"), handles.MachineUpdateHandler)
	machine.DELETE("/delete", middleware.GetAuthMiddleware("admin"), handles.MachineDeleteHandler)
	machine.POST("/add", middleware.GetAuthMiddleware("admin"), handles.MachineAddHandler)

	order := v1.Group("/order") // done except notification
	order.POST("/add", middleware.GetAuthMiddleware("admin"), handles.OrderAddHandler)
	order.PATCH("/update", middleware.GetAuthMiddleware("admin", "worker"), handles.OrderUpdateHandler)
	order.DELETE("/delete", middleware.GetAuthMiddleware("admin"), handles.OrderDeleteHandler)
	order.POST("/query", middleware.GetAuthMiddleware("admin", "worker"), handles.OrderQueryHandler)
	order.GET("/get", middleware.GetAuthMiddleware("admin", "worker"), handles.OrderGetHandler)
	order.GET("/getType", middleware.GetAuthMiddleware("admin"), handles.OrderGetTypeHandler)
	order.PATCH("/updateType", middleware.GetAuthMiddleware("admin"), handles.OrderUpdateTypeHandler)

	version := v1.Group("/version") // done
	version.POST("/query", middleware.GetAuthMiddleware("admin"), handles.VersionQueryHandler)

	notification := v1.Group("/notification") // done
	notification.GET("/list", middleware.GetAuthMiddleware("admin", "worker"), handles.NotificationListHandler)
	notification.POST("/read", middleware.GetAuthMiddleware("admin", "worker"), handles.NotificationReadHandler)
	notification.POST("/readAll", middleware.GetAuthMiddleware("admin", "worker"), handles.NotificationReadAllHandler)

	e.NoRoute(handles.APINoRoute)
}
