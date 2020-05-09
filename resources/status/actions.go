package status

import (
	"backend/resources"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
)

func List(c echo.Context) error {
	var model []Status

	res, err := resources.Fetch(&model, nil, nil)
	if err != nil {
		return err
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, nil)
	} else if res.Error != nil {
		return res.Error
	}

	return c.JSON(http.StatusOK, model)
}
