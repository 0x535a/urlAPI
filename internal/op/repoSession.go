package op

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func newRepo(info *Session) error {
	var err error
	var content []string
	switch info.RepoAPI {
	case "github":
		content, err = util.GetRepo("https://api.github.com/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return errors.WithStack(err)
		}
		util.ListReplacer(&content, "https://raw.githubusercontent.com", database.SettingsStore.Get().Random.SourceRewriteFrom)
	case "gitee":
		content, err = util.GetRepo("https://gitee.com/api/v5/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return errors.WithStack(err)
		}
	default:
		err = errors.WithStack(errors.New(info.RepoAPI + " is not supported"))
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return err
	}
	repoDB := model.Repo{
		UUID:    uuid.New().String(),
		API:     info.RepoAPI,
		Info:    info.RepoInfo,
		Content: string(jsonString),
	}
	return errors.WithStack(db.CreateRepo(&repoDB))
}

func refreshRepo(info *Session) error {
	repoFinder := model.Repo{
		UUID: info.RepoUUID,
	}
	repoDBList, err := db.ReadRepo(repoFinder)
	if err != nil {
		return errors.WithStack(err)
	}
	repoDB := (*repoDBList).RepoList[0]
	info.RepoAPI = repoDB.API
	info.RepoInfo = repoDB.Info
	var content []string
	switch info.RepoAPI {
	case "github":
		content, err = util.GetRepo("https://api.github.com/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return errors.WithStack(err)
		}
		util.ListReplacer(&content, "https://raw.githubusercontent.com", database.SettingsStore.Get().Random.SourceRewriteFrom)
	case "gitee":
		content, err = util.GetRepo("https://gitee.com/api/v5/repos/" + info.RepoInfo + "/contents")
		if err != nil {
			return errors.WithStack(err)
		}
	default:
		err = errors.WithStack(errors.New(info.RepoAPI + " is not supported"))
	}
	jsonString, err := json.Marshal(content)
	if err != nil {
		return errors.WithStack(err)
	}
	repoDB.Content = string(jsonString)
	return errors.WithStack(db.UpdateRepo(&repoDB))
}

func delRepo(info *Session) error {
	repoDB := model.Repo{
		UUID: info.RepoUUID,
	}
	return errors.WithStack(db.DeleteRepo(&repoDB))
}

func fetchRepo(info *Session) error {
	repoFinder := model.Repo{}
	repoDBList, err := db.ReadRepo(repoFinder)
	info.RepoData = repoDBList.RepoList
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return errors.WithStack(err)
	}
	return nil
}
