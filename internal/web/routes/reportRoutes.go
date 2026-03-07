package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/manthan-it-solutions/ams-api/internal/web/controllers"
)

func ReportRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/get-today-report", controllers.GetTodayReport(db))
	router.GET("/get-delivery-report", controllers.GetDeliveryReport(db))
	router.GET("/get-mis-report", controllers.GetMISReport(db))
	router.GET("/get-user-wise-report", controllers.GetUsersWiseReport(db))

	router.GET("/get-today-report-excel", controllers.GetTodayReportExcel(db))
	router.GET("/get-delivery-report-excel", controllers.GetDeliveryReportExcel(db))
	router.GET("/get-mis-report-excel", controllers.GetMISReportExcel(db))
	router.GET("/get-user-wise-report-excel", controllers.GetUsersWiseReportExcel(db))
	router.GET("/get-campaign-wise-report-excel", controllers.GetCampaignWiseReportExcel(db))
}
