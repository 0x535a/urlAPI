package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"urlAPI/internal/model"
)

func (adapter *SQLiteAdapter) CreateRepo(repo *model.Repo) error {
	return errors.WithStack(adapter.db.Create(repo).Error)
}

func (adapter *SQLiteAdapter) UpdateRepo(repo *model.Repo) error {
	if err := adapter.db.Save(repo).Error; err != nil {
		return errors.WithStack(err)
	}
	var tmp []string
	if err := json.Unmarshal([]byte(repo.Content), &tmp); err != nil {
		return errors.WithStack(err)
	}
	RepoMap[repo.API+";"+repo.Info] = tmp
	return nil
}

func (adapter *SQLiteAdapter) ReadRepo(repo model.Repo) (*model.DBList, error) {
	var repos []model.Repo
	var err error
	switch {
	case repo.UUID != "":
		err = adapter.db.Where("uuid = ?", repo.UUID).Find(&repos).Error
	case repo.API != "":
		err = adapter.db.Where("api=? AND info=?", repo.API, repo.Info).Find(&repos).Error
	default:
		err = adapter.db.Find(&repos).Error
	}
	if len(repos) == 0 {
		err = gorm.ErrRecordNotFound
	}
	ret := model.DBList{
		RepoList: repos,
	}
	return &ret, errors.WithStack(err)
}

func (adapter *SQLiteAdapter) DeleteRepo(repo *model.Repo) error {
	delete(RepoMap, repo.API+";"+repo.Info)
	return errors.WithStack(adapter.db.Delete(repo).Error)
}

func CreateRepo(repo *model.Repo) error               { return localDB.CreateRepo(repo) }
func UpdateRepo(repo *model.Repo) error               { return localDB.UpdateRepo(repo) }
func ReadRepo(repo model.Repo) (*model.DBList, error) { return localDB.ReadRepo(repo) }
func DeleteRepo(repo *model.Repo) error               { return localDB.DeleteRepo(repo) }
