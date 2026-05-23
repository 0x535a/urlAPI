package op

import (
	"github.com/pkg/errors"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
)

func GetNotificationByID(notificationID string) (model.Notification, error, int) {
	notification, _, err := db.GetNotificationByID(notificationID)
	if err != nil {
		return model.Notification{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return notification, nil, 0
}

func GetNotificationByUserID(userID string) ([]model.Notification, error, int) {
	notifications, _, err := db.GetNotificationByUserID(userID)
	if err != nil {
		return nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return notifications, nil, 0
}

func DeleteNotificationByID(notificationID string) (error, int) {
	if err := db.DeleteNotificationByID(notificationID); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func DeleteNotificationsByUserID(userID string) (error, int) {
	if err := db.DeleteNotificationByUserID(userID); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}
