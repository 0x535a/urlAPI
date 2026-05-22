package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"urlAPI/internal/database"
	"urlAPI/util"

	"github.com/gin-gonic/gin"
)

const skipDBKey = "security.skipDB"
const generalKey = "security.general"

func GeneralSecurityMiddleware(kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		general := General{
			Referer: c.Request.Referer(),
			IP:      c.ClientIP(),
			Time:    time.Now(),
		}
		switch kind {
		case "txt":
			general.Type = util.TypeMap["txt.gen"]
			general.Target = c.Query("prompt")
			general.Info = c.Query("more")
		case "img":
			general.Type = util.TypeMap["img.gen"]
			general.Target = c.Query("prompt")
			general.Info = c.Query("more")
		case "rand":
			general.Type = util.TypeMap["rand"]
			general.Target = c.Query("user") + "/" + c.Query("repo")
			general.Info = c.Query("more")
		case "web":
			general.Type = util.TypeMap["web.img"]
			general.Target = c.Query("img")
			general.Info = c.Query("more")
		case "download":
			general.Type = util.TypeMap["download"]
			general.Target = c.Query("img")
			general.SkipDB = true
		}
		checkGeneralSecurity(&general)
		if general.Unsafe {
			log.Printf("%s from %s\n", general.Info, c.ClientIP())
			c.JSON(http.StatusForbidden, gin.H{"error": general.Info})
			c.Abort()
			return
		}
		c.Set(skipDBKey, general.SkipDB)
		c.Set(generalKey, general)
		c.Next()
	}
}

func TextSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		general, ok := GetGeneral(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "security context missing"})
			return
		}
		if !general.Unsafe {
			settings := database.SettingsStore.Get()
			txtgenenabled := settings.Text.EnabledPromptKeys
			txtacceptprompt := settings.Text.AcceptedPromptGlob
			textAPIs := []string{"openai", "alibaba", "deepseek", "otherapi"}
			api := c.Query("api")
			var prompt string
			if _, ok := database.PromptMap[general.Target]; ok {
				prompt = general.Target
			} else {
				prompt = "other"
			}
			switch {
			case !settings.Features.TextEnabled:
				general.Info = "Txt is not enabled"
			case !util.ListChecker(&txtgenenabled, &prompt):
				general.Info = fmt.Sprintf("Target %s is not enabled for Txt Gen", general.Target)
			case (general.Target == "" || !util.WildcardChecker(&txtacceptprompt, &general.Target)) && prompt == "other":
				general.Info = fmt.Sprintf("Prompt %s is not enabled for Txt Gen", general.Target)
			case !util.ListChecker(&textAPIs, &api):
				general.Info = "Invalid API"
			}
			if general.Info != "" {
				general.Unsafe = true
			}
		}
		abortIfUnsafe(c, general)
	}
}

func ImageSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		general, ok := GetGeneral(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "security context missing"})
			return
		}
		if !general.Unsafe {
			settings := database.SettingsStore.Get()
			imgacceptprompt := settings.Image.AcceptedPromptGlob
			imgAPIs := []string{"openai", "alibaba"}
			api := c.Query("api")
			if !settings.Features.ImageEnabled {
				general.Info = "Img is not enabled"
			} else if general.Target == "" || !util.WildcardChecker(&imgacceptprompt, &general.Target) {
				general.Info = fmt.Sprintf("Prompt %s is not allowed for ImgGen", general.Target)
			} else if !util.ListChecker(&imgAPIs, &api) {
				general.Info = "Invalid API"
			}
			if general.Info != "" {
				general.Unsafe = true
			}
		}
		abortIfUnsafe(c, general)
	}
}

func RandomSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		general, ok := GetGeneral(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "security context missing"})
			return
		}
		if !general.Unsafe {
			settings := database.SettingsStore.Get()
			randAPIs := []string{"github", "gitee"}
			api := c.Query("api")
			if !settings.Features.RandomEnabled {
				general.Info = "Random is not enabled"
			} else if !util.ListChecker(&randAPIs, &api) {
				general.Info = "Invalid API"
			}
			if general.Info != "" {
				general.Unsafe = true
			}
		}
		abortIfUnsafe(c, general)
	}
}

func WebSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		general, ok := GetGeneral(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "security context missing"})
			return
		}
		parsedURL, _ := url.Parse(general.Target)
		if !general.Unsafe {
			settings := database.SettingsStore.Get()
			webimgallowed := settings.Web.AllowedHosts
			api := parsedURL.Host
			switch {
			case !settings.Features.WebImgEnabled:
				general.Info = "WebImg is not enabled"
			case !util.ListChecker(&webimgallowed, &api):
				general.Info = fmt.Sprintf("API %s is not enabled", api)
			case api == "www.ithome.com" && !settings.Features.TextEnabled:
				general.Info = "For ITHome, TxtSum is not enabled"
			}
			if general.Info != "" {
				general.Unsafe = true
			}
		}
		abortIfUnsafe(c, general)
	}
}

func DownloadSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func GetGeneral(c *gin.Context) (General, bool) {
	value, exists := c.Get(generalKey)
	if !exists {
		return General{}, false
	}
	general, ok := value.(General)
	return general, ok
}

func GetSkipDB(c *gin.Context) bool {
	value, exists := c.Get(skipDBKey)
	if !exists {
		return false
	}
	skipDB, ok := value.(bool)
	return ok && skipDB
}

func abortIfUnsafe(c *gin.Context, general General) {
	if general.Unsafe {
		log.Printf("%s from %s\n", general.Info, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{"error": general.Info})
		c.Abort()
		return
	}
	c.Set(skipDBKey, general.SkipDB)
	c.Set(generalKey, general)
	c.Next()
}
