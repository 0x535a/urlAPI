package database

import (
	"gorm.io/gorm/clause"
	"zhongxin/internal/model"
	"zhongxin/util"
)

func (s *SQLiteAdapter) NewToken(id, token, typ string, expireSeconds int64) error {
	newToken := model.Token{
		UserID:     id,
		Token:      token,
		ExpireTime: expireSeconds,
		UserType:   typ,
	}
	return s.db.Create(&newToken).Error
}

func (s *SQLiteAdapter) GetToken(token string) (model.Token, bool, error) {
	var t model.Token
	result := s.db.Where(clause.Eq{Column: "token", Value: token}).First(&t)
	if result.Error != nil {
		return t, false, result.Error
	}
	if result.RowsAffected == 0 {
		return t, false, nil
	}
	return t, true, nil
}

func (s *SQLiteAdapter) CleanExpiredTokens() error {
	now := util.TimeNow().Unix()
	return s.db.Where("expireTime <= ?", now).Delete(&model.Token{}).Error
}
