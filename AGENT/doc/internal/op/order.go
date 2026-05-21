package op

import (
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewOrder(order model.Order) (error, int) {
	// Attention: haven't generated the name of the order
	if err := db.CreateOrder(order); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func GetOrderByID(orderID string) (model.Order, error, int) {
	order, _, err := db.GetOrderByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Order{}, errors.WithStack(errors.New("order not found")), _error.DBRecordNotFound
		}
		return model.Order{}, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return order, nil, 0
}

func UpdateOrder(order model.Order, userID string) (error, int) {
	if userID == "" {
		err := db.UpdateOrderWithoutVersion(order)
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		return nil, 0
	}

	oldVersionID, newVersionID, err := db.UpdateOrder(order)
	if err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	if err := db.NewVersion(userID, order.ID, oldVersionID, newVersionID, "order"); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func DeleteOrderByID(orderID string) (error, int) {
	if err := db.DeleteOrderByID(orderID); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	if err := db.DeleteVersionByContentID(orderID); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}

func GetOrdersByFilter(keys, rels []string, vals []interface{}) ([]model.Order, error, int) {
	cls := make([]clause.Expression, len(keys))
	for index, key := range keys {
		if !util.IsComparable(rels[index]) {
			continue
		}
		switch rels[index] {
		case "ge":
			cls[index] = clause.Gte{Column: key, Value: vals[index]}
		case "le":
			cls[index] = clause.Lte{Column: key, Value: vals[index]}
		case "eq":
			cls[index] = clause.Eq{Column: key, Value: vals[index]}
		case "has":
			if s, ok := vals[index].(string); ok {
				cls[index] = clause.Like{Column: key, Value: "%" + s + "%"}
			}
		}
	}
	orders, _, err := db.GetOrdersByFilter(cls)
	if err != nil {
		return nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return orders, nil, 0
}

func GetOrderTypes() ([]string, error, int) {
	types, ok, err := db.GetKVByKey("orderType")
	if !ok {
		return []string{}, nil, 0
	}
	if err != nil {
		return nil, errors.WithStack(err), _error.ConvertGormError(err)
	}
	return util.DBToStringList(types.Value), nil, 0
}

func UpdateOrderTypes(types []string) (error, int) {
	if err := db.SetKV("orderType", util.StringListToDB(types)); err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	return nil, 0
}
