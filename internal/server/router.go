package server

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"urlAPI/internal/server/handles"
	"urlAPI/internal/server/middleware"
	"urlAPI/static"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST"}

	r := gin.New()
	r.Use(cors.New(config))
	rootFS, _ := fs.Sub(static.StaticFS, "dist")
	assetsFS, _ := fs.Sub(static.StaticFS, "dist/assets")
	r.StaticFS("/assets", http.FS(assetsFS))
	tpl := template.Must(template.ParseFS(rootFS, "*.html"))
	r.SetHTMLTemplate(tpl)

	r.GET("/txt", middleware.GeneralSecurityMiddleware("txt"), middleware.TextSecurityMiddleware(), handles.TxtHandler)
	r.GET("/img", middleware.GeneralSecurityMiddleware("img"), middleware.ImageSecurityMiddleware(), handles.ImgHandler)
	r.GET("/rand", middleware.GeneralSecurityMiddleware("rand"), middleware.RandomSecurityMiddleware(), handles.RandHandler)
	r.GET("/web", middleware.GeneralSecurityMiddleware("web"), middleware.WebSecurityMiddleware(), handles.WebHandler)
	r.GET("/download", middleware.GeneralSecurityMiddleware("download"), middleware.DownloadSecurityMiddleware(), handles.DownloadHandler)
	r.POST("/session", handles.SessionHandler)
	r.NoRoute(handles.StaticHandler)
	return r
}

func Run(port string) error {
	log.Printf("The server will be running on port %s", port)
	return NewRouter().Run(":" + port)
}
