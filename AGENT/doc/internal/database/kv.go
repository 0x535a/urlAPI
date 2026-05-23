package database

import (
	"zhongxin/internal/model"

	"gorm.io/gorm/clause"
)

func (s *SQLiteAdapter) SetKV(key string, val string) error {
	kv := model.KV{
		Key:   key,
		Value: val,
	}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&kv).Error
}

func (s *SQLiteAdapter) GetKVByKey(key string) (model.KV, bool, error) {
	var kv model.KV
	result := s.db.Where("key = ?", key).First(&kv)
	if result.Error != nil {
		return kv, false, result.Error
	}
	if result.RowsAffected == 0 {
		return kv, false, nil
	}
	return kv, true, nil
}

func (s *SQLiteAdapter) DeleteKVByKey(key string) error {
	return s.db.Where("key = ?", key).Delete(&model.KV{}).Error
}

func (s *SQLiteAdapter) GetAllKVs() ([]model.KV, error) {
	var kvs []model.KV
	result := s.db.Find(&kvs)
	if result.Error != nil {
		return kvs, result.Error
	}
	return kvs, nil
}

func (s *SQLiteAdapter) GetKVsByFilter(cls []clause.Expression) ([]model.KV, error) {
	var kvs []model.KV
	result := s.db
	for _, c := range cls {
		result = result.Where(c)
	}
	result = result.Find(&kvs)
	if result.Error != nil {
		return nil, result.Error
	}
	return kvs, nil
}
