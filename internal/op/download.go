package op

import (
	"os"
	"urlAPI/file"
	"urlAPI/internal/database"
)

func downloadImage(target string) ([]byte, string, error) {
	var img []byte
	var err error
	switch target {
	case "empty":
		img, err = file.EmptyPNG.ReadFile("empty.png")
	default:
		img, err = os.ReadFile(ImgPath + target + ".png")
	}
	if err != nil {
		return nil, database.SettingsStore.Get().Web.FallbackImageURL, err
	}
	return img, "", nil
}
