package database

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"urlAPI/internal/model"
)

func (adapter *SQLiteAdapter) CreateSession(session *model.Session) error {
	SessionMap[session.Token] = *session
	return errors.WithStack(adapter.db.Create(&session).Error)

}

func (adapter *SQLiteAdapter) UpdateSession(session *model.Session) error {
	SessionMap[session.Token] = *session
	return errors.WithStack(adapter.db.Save(session).Error)
}

func (adapter *SQLiteAdapter) ReadSession(session model.Session) (*model.DBList, error) {
	var sessions []model.Session
	err := adapter.db.Where("token=?", session.Token).Find(&sessions).Error
	if len(sessions) == 0 {
		err = gorm.ErrRecordNotFound
	}
	ret := model.DBList{
		SessionList: sessions,
	}
	return &ret, errors.WithStack(err)
}

func (adapter *SQLiteAdapter) DeleteSession(session *model.Session) error {
	delete(SessionMap, session.Token)
	return errors.WithStack(adapter.db.Delete(session).Error)
}

func CreateSession(session *model.Session) error               { return localDB.CreateSession(session) }
func UpdateSession(session *model.Session) error               { return localDB.UpdateSession(session) }
func ReadSession(session model.Session) (*model.DBList, error) { return localDB.ReadSession(session) }
func DeleteSession(session *model.Session) error               { return localDB.DeleteSession(session) }
