package op

import (
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
)

func GetVersionByFilter(cls []clause.Expression) ([]model.Version, error, int) {
	ret, err := db.GetVersionByFilter(cls)
	if err != nil {
		return nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return ret, nil, 0
}
