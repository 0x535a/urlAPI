package database

import (
	"gorm.io/gorm/clause"
	"zhongxin/internal/model"
)

func (s *SQLiteAdapter) NewUser(user model.User) error {
	return s.db.Create(&user).Error
}

func (s *SQLiteAdapter) DeleteUserByID(id string) error {
	return s.db.Where(clause.Eq{Column: "id", Value: id}).Delete(&model.User{}).Error
}

func (s *SQLiteAdapter) DeleteUserByName(name string) error {
	return s.db.Where(clause.Eq{Column: "name", Value: name}).Delete(&model.User{}).Error
}

func (s *SQLiteAdapter) UpdateUser(user model.User) error {
	return s.db.Save(&user).Error
}

func (s *SQLiteAdapter) GetUserByID(id string) (model.User, bool, error) {
	var u model.User
	result := s.db.Where(clause.Eq{Column: "id", Value: id}).First(&u)
	if result.Error != nil {
		return u, false, result.Error
	}
	if result.RowsAffected == 0 {
		return u, false, nil
	}
	return u, true, nil
}

func (s *SQLiteAdapter) GetUserByWXID(wxid string) (model.User, bool, error) {
	var u model.User
	result := s.db.Where(clause.Eq{Column: "wxid", Value: wxid}).First(&u)
	if result.Error != nil {
		return u, false, result.Error
	}
	if result.RowsAffected == 0 {
		return u, false, nil
	}
	return u, true, nil
}

func (s *SQLiteAdapter) GetUserByName(name string) (model.User, bool, error) {
	var u model.User
	result := s.db.Where(clause.Eq{Column: "name", Value: name}).First(&u)
	if result.Error != nil {
		return u, false, result.Error
	}
	if result.RowsAffected == 0 {
		return u, false, nil
	}
	return u, true, nil
}

func (s *SQLiteAdapter) GetUserByFilter(cls []clause.Eq) ([]model.User, bool, error) {
	var users []model.User
	result := s.db
	for _, cl := range cls {
		result = result.Where(cl)
	}
	result.Find(&users)
	if result.Error != nil {
		return users, false, result.Error
	}
	if result.RowsAffected == 0 {
		return users, false, nil
	}
	return users, true, nil
}

func (s *SQLiteAdapter) GetUsersByType(userType string) ([]model.User, bool, error) {
	var us []model.User
	result := s.db.Where(clause.Eq{Column: "type", Value: userType}).Find(&us)
	if result.Error != nil {
		return us, false, result.Error
	}
	if result.RowsAffected == 0 {
		return us, false, nil
	}
	return us, true, nil
}
