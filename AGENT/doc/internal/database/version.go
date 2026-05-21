package database

import (
	"gorm.io/gorm/clause"
	"zhongxin/internal/model"
	"zhongxin/util"
)

func (s *SQLiteAdapter) NewVersion(userID string, contentID string, from string, to string, vType string) error {
	versionID := util.NewUUID()
	version := model.Version{
		ID:        versionID,
		Type:      vType,
		UserID:    userID,
		ContentID: contentID,
		From:      from,
		To:        to,
		Time:      util.TimeNow().Unix(),
	}
	if err := s.db.Create(&version).Error; err != nil {
		return err
	}
	return nil
}

func (s *SQLiteAdapter) GetVersionByFilter(cls []clause.Expression) ([]model.Version, error) {
	var versions []model.Version
	result := s.db
	for _, c := range cls {
		result = result.Where(c)
	}
	result = result.Find(&versions)
	if result.Error != nil {
		return nil, result.Error
	}
	return versions, nil
}

func (s *SQLiteAdapter) DeleteVersionByContentID(id string) error {
	return s.db.Where(clause.Eq{Column: "contentID", Value: id}).Delete(&model.Version{}).Error
}
