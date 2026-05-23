package op

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func getBiliABV(URL string) string {
	for i := 31; i < len(URL); i++ {
		if URL[i] == '/' || URL[i] == '?' {
			return URL[31:i]
		}
	}
	return URL[31:]
}

func getYtbID(URL string) string {
	for i := 32; i < len(URL); i++ {
		if URL[i] == '&' {
			return URL[32:i]
		}
	}
	return URL[32:]
}

func generateWebImage(task *model.Task, host string) (GenerateResult, error) {
	var img []byte
	var err error
	settings := database.SettingsStore.Get()
	switch task.API {
	case "www.bilibili.com":
		img, err = util.Bili(getBiliABV(task.Target))
	case "www.youtube.com":
		img, err = util.Ytb(getYtbID(task.Target), settings.Web.YouTubeToken)
	case "arxiv.org":
		img, err = util.Arxiv(task.Target)
	case "www.ithome.com":
		api := settings.Text.SummaryAPI
		provider, ok := settings.Providers.ByName(api)
		if !ok {
			task.Status = "failed"
			task.Return = "Invalid API"
			return GenerateResult{}, fmt.Errorf("web image invalid summary API")
		}
		img, err = util.ITHome(task.Target, provider.Endpoint, provider.APIKey, provider.SummaryModel, settings.Prompts.SummaryContext)
	case "github.com", "gitee.com":
		img, err = util.Repo(task.Target, settings.Web.RepoToken)
	default:
		task.Status = "failed"
		task.Return = "Invalid URL"
		return GenerateResult{}, fmt.Errorf("web image invalid URL")
	}
	if err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return GenerateResult{}, fmt.Errorf("web image generation: %w", err)
	}
	file, err := os.Create(ImgPath + task.UUID + ".png")
	if err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return GenerateResult{}, fmt.Errorf("web image create: %w", err)
	}
	defer file.Close()
	if _, err = io.Copy(file, bytes.NewReader(img)); err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return GenerateResult{}, fmt.Errorf("web image write: %w", err)
	}
	result := GenerateResult{URL: host + "/download?img=" + task.UUID}
	if err := setTaskResult(task, result); err != nil {
		return GenerateResult{}, err
	}
	return result, nil
}
