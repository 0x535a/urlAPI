package error

import (
	"errors"
	"gorm.io/gorm"
)

// 自定义错误码定义
const (
	DBOtherError = iota + 300
	DBRecordNotFound
	DBInvalidTransaction
	DBNotImplemented
	DBMissingWhereClause
	DBUnsupportedRelation
	DBPrimaryKeyRequired
	DBModelValueRequired
	DBModelAccessibleFieldsRequired
	DBSubQueryRequired
	DBInvalidData
	DBUnsupportedDriver
	DBRegistered
	DBInvalidField
	DBEmptySlice
	DBDryRunModeUnsupported
	DBInvalidDB
	DBInvalidValue
	DBInvalidValueOfLength
	DBPreloadNotAllowed
	DBDuplicatedKey
	DBForeignKeyViolated
	DBCheckConstraintViolated
)

// ConvertGormError 将gorm错误转换为自定义错误码
func ConvertGormError(err error) int {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return DBRecordNotFound
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return DBInvalidTransaction
	case errors.Is(err, gorm.ErrNotImplemented):
		return DBNotImplemented
	case errors.Is(err, gorm.ErrMissingWhereClause):
		return DBMissingWhereClause
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return DBUnsupportedRelation
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return DBPrimaryKeyRequired
	case errors.Is(err, gorm.ErrModelValueRequired):
		return DBModelValueRequired
	case errors.Is(err, gorm.ErrModelAccessibleFieldsRequired):
		return DBModelAccessibleFieldsRequired
	case errors.Is(err, gorm.ErrSubQueryRequired):
		return DBSubQueryRequired
	case errors.Is(err, gorm.ErrInvalidData):
		return DBInvalidData
	case errors.Is(err, gorm.ErrUnsupportedDriver):
		return DBUnsupportedDriver
	case errors.Is(err, gorm.ErrRegistered):
		return DBRegistered
	case errors.Is(err, gorm.ErrInvalidField):
		return DBInvalidField
	case errors.Is(err, gorm.ErrEmptySlice):
		return DBEmptySlice
	case errors.Is(err, gorm.ErrDryRunModeUnsupported):
		return DBDryRunModeUnsupported
	case errors.Is(err, gorm.ErrInvalidDB):
		return DBInvalidDB
	case errors.Is(err, gorm.ErrInvalidValue):
		return DBInvalidValue
	case errors.Is(err, gorm.ErrInvalidValueOfLength):
		return DBInvalidValueOfLength
	case errors.Is(err, gorm.ErrPreloadNotAllowed):
		return DBPreloadNotAllowed
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return DBDuplicatedKey
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return DBForeignKeyViolated
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		return DBCheckConstraintViolated
	default:
		return DBOtherError
	}
}
