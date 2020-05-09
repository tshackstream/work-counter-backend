package projects

import (
	"backend/resources"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
)

func Get(c echo.Context) error {
	model := Project{}
	condition := []resources.WhereCondition{
		{Column: "id", Operation: "=", Value: c.Param("project_id")},
	}

	res, err := resources.FetchOne(&model, condition)
	if err != nil {
		return err
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusOK, nil)
	} else if res.Error != nil {
		return res.Error
	}

	return c.JSON(http.StatusOK, model)
}

func CreateOrUpdate(c echo.Context) error {
	model := new(Project)

	if err := c.Bind(model); err != nil {
		return err
	}

	res, err := resources.Save(&model)

	if err != nil {
		return err
	}

	if res.Error != nil {
		return res.Error
	}

	return c.JSON(http.StatusOK, model)
}
