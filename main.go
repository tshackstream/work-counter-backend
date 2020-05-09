package main

import (
	"backend/resources/download"
	"backend/resources/monthly_work_results"
	"backend/resources/projects"
	"backend/resources/status"
	"backend/resources/work_info"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/leekchan/timeutil"
	"net/http"
	"os"
	"time"
)

func main() {
	e := echo.New()
	t := time.Now()
	fp, err := os.OpenFile("logs/"+timeutil.Strftime(&t, "%Y-%m-%d")+".log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		panic(err)
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: fp,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("ALLOW_ORIGIN")},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	}))

	routing(e)
	port := os.Getenv("ECHO_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

func routing(e *echo.Echo) {
	e.GET("/api/v1/projects/:project_id", projects.Get)
	e.POST("/api/v1/projects", projects.CreateOrUpdate)
	e.PUT("/api/v1/projects", projects.CreateOrUpdate)
	e.GET("/api/v1/work_info/list/:project_id/:year/:month", work_info.List)
	e.PUT("/api/v1/work_info", work_info.BulkUpdate)
	e.GET("/api/v1/status/list", status.List)
	e.GET("/api/v1/monthly_work_result/:project_id/:year/:month", monthly_work_results.Get)
	e.POST("/api/v1/monthly_work_result", monthly_work_results.CreateOrUpdate)
	e.GET("/api/v1/download/work_sheet/:project_id/:year/:month", download.WorkSheet)
	e.GET("/api/v1/download/invoice/:project_id/:year/:month", download.Invoice)
}
