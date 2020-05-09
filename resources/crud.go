package resources

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type WhereCondition struct {
	Column    string
	Operation string
	Value     interface{}
}

type Order struct {
	Column    string
	Operation interface{}
}

func FetchOne(model interface{}, conditions []WhereCondition) (*gorm.DB, error) {
	db, err := DbConnect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if conditions != nil {
		db, err = setCondition(db, conditions)
		if err != nil {
			return nil, err
		}
	}

	return db.First(model), nil
}

func Fetch(model interface{}, conditions []WhereCondition, order []Order) (*gorm.DB, error) {
	db, err := DbConnect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if conditions != nil {
		db, err = setCondition(db, conditions)
		if err != nil {
			return nil, err
		}
	}

	if order != nil {
		db, err = setOrder(db, order)
		if err != nil {
			return nil, err
		}
	}

	return db.Find(model), nil
}

func Save(param interface{}) (*gorm.DB, error) {
	db, err := DbConnect()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	res := db.Save(param)

	return res, nil
}

func setCondition(db *gorm.DB, conditions []WhereCondition) (*gorm.DB, error) {
	for _, condition := range conditions {
		switch condition.Operation {
		case "BETWEEN":
			if value, ok := condition.Value.([]string); ok {
				db = db.Where(fmt.Sprintf("%s %s ? AND ?", condition.Column, condition.Operation), value[0], value[1])
			} else {
				return nil, errors.New("invalid Where BETWEEN condition value")
			}
		default:
			db = db.Where(fmt.Sprintf("%s %s ?", condition.Column, condition.Operation), condition.Value)
		}
	}

	return db, nil
}

func setOrder(db *gorm.DB, orders []Order) (*gorm.DB, error) {
	for _, order := range orders {
		var ope string
		if order.Operation == nil {
			ope = "asc"
		} else {
			ope = order.Operation.(string)
		}

		query := order.Column + " " + ope
		db.Order(query)
	}

	return db, nil
}
