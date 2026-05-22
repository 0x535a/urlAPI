package handles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"urlAPI/internal/op"
	"urlAPI/util"
)

func DownloadHandler(c *gin.Context) {
	img, fallbackURL, err := op.DownloadImage(c.Query("img"))
	util.ErrorPrinter(err)
	downloadReturn(c, img, fallbackURL)
}

func downloadReturn(c *gin.Context, img []byte, fallbackURL string) {
	if fallbackURL != "" {
		c.Redirect(http.StatusFound, fallbackURL)
		return
	}
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="download.png"`))
	c.Header("Accept-Length", fmt.Sprintf("%d", len(img)))
	c.Writer.Write(img)
}
