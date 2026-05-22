package security

import (
	"fmt"
	"urlAPI/database"
	"urlAPI/util"
)

func (info *TxtGen) FunctionChecker(general *General) {
	settings := database.SettingsStore.Get()
	txtgenenabled := settings.Text.EnabledPromptKeys
	txtacceptprompt := settings.Text.AcceptedPromptGlob
	var prompt string
	if _, ok := database.PromptMap[general.Target]; ok {
		prompt = general.Target
	} else {
		prompt = "other"
	}
	switch {
	case !settings.Features.TextEnabled:
		general.Info = "Txt is not enabled"
		break
	case !util.ListChecker(&txtgenenabled, &(prompt)):
		general.Info = fmt.Sprintf("Target %s is not enabled for Txt	Gen", general.Target)
	case (general.Target == "" || !util.WildcardChecker(&txtacceptprompt, &(general.Target))) && prompt == "other":
		general.Info = fmt.Sprintf("Prompt %s is not enabled for Txt	Gen", general.Target)
		break
	default:
		return
	}
	general.Unsafe = true
}

func (info *TxtSum) FunctionChecker(general *General) {
	settings := database.SettingsStore.Get()
	switch {
	case !settings.Features.TextEnabled:
		general.Info = "Txt is not enabled"
		break
	default:
		return
	}
	general.Unsafe = true
}

func (info *ImgGen) FunctionChecker(general *General) {
	settings := database.SettingsStore.Get()
	imgacceptprompt := settings.Image.AcceptedPromptGlob
	switch {
	case !settings.Features.ImageEnabled:
		general.Info = "Img is not enabled"
		break
	case general.Target == "" || !util.WildcardChecker(&imgacceptprompt, &(general.Target)):
		general.Info = fmt.Sprintf("Prompt %s is not allowed for ImgGen", general.Target)
	default:
		return
	}
	general.Unsafe = true
}

func (info *Rand) FunctionChecker(general *General) {
	settings := database.SettingsStore.Get()
	switch {
	case !settings.Features.RandomEnabled:
		general.Info = "Random is not enabled"
		break
	default:
		return
	}
	general.Unsafe = true
}

func (info *WebImg) FunctionChecker(general *General) {
	settings := database.SettingsStore.Get()
	webimgallowed := settings.Web.AllowedHosts

	switch {
	case !settings.Features.WebImgEnabled:
		general.Info = "WebImg is not enabled"
		break
	case !util.ListChecker(&webimgallowed, &(info.API)):
		general.Info = fmt.Sprintf("API %s is not enabled", info.API)
		break
	case info.API == "www.ithome.com" && !settings.Features.TextEnabled:
		general.Info = "For ITHome, TxtSum is not enabled"
		break
	default:
		return
	}
	general.Unsafe = true
}
