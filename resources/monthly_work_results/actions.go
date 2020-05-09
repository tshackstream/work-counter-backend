package monthly_work_results

import (
	"backend/resources"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
)

func Get(c echo.Context) error {
	model := MonthlyWorkResult{}
	condition := []resources.WhereCondition{
		{Column: "project_id", Operation: "=", Value: c.Param("project_id")},
		{Column: "month", Operation: "=", Value: c.Param("year") + "-" + c.Param("month") + "-01"},
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
	model := new(MonthlyWorkResult)

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
