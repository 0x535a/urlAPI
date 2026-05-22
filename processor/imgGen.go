package processor

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"urlAPI/database"
	"urlAPI/util"
)

func (info *ImgGen) Process(data *database.Task) error {
	settings := database.SettingsStore.Get()
	info.Return = settings.Image.FallbackImageURL
	if info.API == "" {
		info.API = settings.Image.API
		data.API = info.API
	}
	provider, ok := settings.Providers.ByName(info.API)
	if !ok {
		data.Status = "failed"
		data.Return = "Imggen Process invalid API"
		return errors.WithStack(errors.New("Imggen Process invalid API"))
	}
	if info.Model == "" {
		info.Model = provider.ImageModel
		data.Model = info.Model
	}
	if info.Size == "" {
		info.Size = provider.ImageSize
		data.Size = info.Size
	}
	token := provider.APIKey
	var img []byte
	prompt := info.Target
	var err error
	switch info.API {
	case "alibaba":
		img, prompt, err = util.AlibabaImg(token, info.Target, info.Model, info.Size)
	case "openai":
		prompt = info.Target
		img, err = util.OpenaiImg(settings.Providers.OpenAI.Endpoint,
			token, info.Target, info.Model, info.Size)
	default:
		data.Status = "failed"
		data.Return = "Imggen Process invalid API"
		return errors.WithStack(errors.New("Imggen Process invalid API"))
	}
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.WithStack(err)
	}
	file, err := os.Create(ImgPath + data.UUID + ".png")
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.WithStack(err)
	}
	_, err = io.Copy(file, bytes.NewReader(img))
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.WithStack(err)
	}
	data.Return = fmt.Sprintf(`{"original_prompt": "%s", "actual_prompt": "%s", "url": "%s"}`, info.Target, prompt, info.Host+"/download?img="+data.UUID)
	data.Status = "success"
	info.Return = info.Host + "/download?img=" + data.UUID
	defer file.Close()
	return nil
}
