package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manthan-it-solutions/ams-api/internal/web/models"
	"github.com/manthan-it-solutions/ams-api/internal/web/utils"
	"github.com/xuri/excelize/v2"
)

func GetTodayReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetRequestParameters(c)

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":
		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)
		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, `
				(campaign_title LIKE ? OR 
				 sender_mobile LIKE ? OR 
				 voice_status LIKE ?)
			`)
			params = append(params, searchTerm, searchTerm, searchTerm)
		}

		conditions = append(conditions, "send_date = CURDATE()")

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		countQuery := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM temp_obd_send_tbl_final
			%s
		`, whereClause)

		var totalCount int
		err := db.QueryRow(countQuery, params...).Scan(&totalCount)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to count records")
			return
		}

		query := fmt.Sprintf(`
			SELECT 
				campaign_title, 
				campaign_no,
				sender_mobile, 
				temp_category,
				voice_status, 
				send_time
			FROM temp_obd_send_tbl_final
			%s
			ORDER BY send_time DESC
			LIMIT ? OFFSET ?
		`, whereClause)

		queryParams := append([]interface{}{}, params...)
		queryParams = append(queryParams, limit, offset)

		rows, err := db.Query(query, queryParams...)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to fetch records")
			return
		}
		defer rows.Close()

		var results []models.GetTodayReportPayload

		for rows.Next() {
			var item models.GetTodayReportPayload
			err := rows.Scan(
				&item.CampaignTitle,
				&item.CampaignNo,
				&item.SenderMobile,
				&item.TempCategory,
				&item.VoiceStatus,
				&item.SendTime,
			)
			if err != nil {
				utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to scan record")
				return
			}
			results = append(results, item)
		}

		c.JSON(http.StatusOK, gin.H{
			"total": totalCount,
			"page":  page,
			"limit": limit,
			"data":  results,
		})
	}
}

func GetDeliveryReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetRequestParameters(c)

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		fromDate := c.DefaultQuery("fromDate", "")
		toDate := c.DefaultQuery("toDate", "")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":
		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)
		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if fromDate != "" {
			conditions = append(conditions, "send_date >= ?")
			params = append(params, fromDate)
		}

		if toDate != "" {
			conditions = append(conditions, "send_date <= ?")
			params = append(params, toDate)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, `
				(campaign_title LIKE ? OR 
				 sender_mobile LIKE ? OR 
				 voice_status LIKE ?)
			`)
			params = append(params, searchTerm, searchTerm, searchTerm)
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}

		countQuery := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM temp_obd_send_tbl_final
			%s
		`, whereClause)

		var totalCount int
		err := db.QueryRow(countQuery, params...).Scan(&totalCount)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to count records")
			return
		}

		dataQuery := fmt.Sprintf(`
			SELECT 
				campaign_title, 
				campaign_no,
				business_no,
				sender_mobile, 
				temp_name,
				temp_category,
				voice_status, 
				COALESCE(call_duration," ") AS call_duration,
				send_date,
				send_time,
				count_voice_deduct
			FROM temp_obd_send_tbl_final
			%s
			ORDER BY send_date DESC, send_time DESC
			LIMIT ? OFFSET ?
		`, whereClause)

		queryParams := append([]interface{}{}, params...)
		queryParams = append(queryParams, limit, offset)

		rows, err := db.Query(dataQuery, queryParams...)
		if err != nil {
			utils.LogError(c, err)
			return
		}
		defer rows.Close()

		var results []models.GetDeliveryReportPayload

		for rows.Next() {
			var item models.GetDeliveryReportPayload
			err := rows.Scan(
				&item.CampaignTitle,
				&item.CampaignNo,
				&item.BusinessNo,
				&item.SenderMobile,
				&item.TempName,
				&item.TempCategory,
				&item.VoiceStatus,
				&item.CallDuration,
				&item.SendDate,
				&item.SendTime,
				&item.CountVoiceDeduct,
			)
			if err != nil {
				utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to scan record")
				return
			}
			results = append(results, item)
		}

		c.JSON(http.StatusOK, gin.H{
			"total": totalCount,
			"page":  page,
			"limit": limit,
			"data":  results,
		})
	}
}

func GetMISReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetRequestParameters(c)

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		fromDate := c.DefaultQuery("fromDate", "")
		toDate := c.DefaultQuery("toDate", "")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":
		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)
		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if fromDate != "" {
			conditions = append(conditions, "report_date >= ?")
			params = append(params, fromDate)
		}

		if toDate != "" {
			conditions = append(conditions, "report_date <= ?")
			params = append(params, toDate)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, "(report_date LIKE ? OR answered LIKE ?)")
			params = append(params, searchTerm, searchTerm)
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}
		countQuery := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM mis_report_tbl
			%s
		`, whereClause)

		var totalCount int
		err := db.QueryRow(countQuery, params...).Scan(&totalCount)
		if err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count records"})
			return
		}

		dataQuery := fmt.Sprintf(`
			SELECT 
			    user_id,
				DATE_FORMAT(report_date,'%%d-%%m-%%Y') as report_date,
				answered,
				not_answered,
				total_calls
			FROM mis_report_tbl
			%s
			ORDER BY report_date DESC
			LIMIT ? OFFSET ?
		`, whereClause)

		queryParams := make([]interface{}, len(params))
		copy(queryParams, params)
		queryParams = append(queryParams, limit, offset)

		rows, err := db.Query(dataQuery, queryParams...)
		if err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
			return
		}
		defer rows.Close()

		var results []models.GetMISReportPayload

		for rows.Next() {
			var item models.GetMISReportPayload
			err := rows.Scan(
				&item.UserId,
				&item.ReportDate,
				&item.Answered,
				&item.NotAnswered,
				&item.TotalsCalls,
			)
			if err != nil {
				utils.LogError(c, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan record"})
				return
			}
			results = append(results, item)
		}

		if err = rows.Err(); err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error iterating rows"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total": totalCount,
			"page":  page,
			"limit": limit,
			"data":  results,
		})
	}
}

func GetUsersWiseReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetRequestParameters(c)

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		fromDate := c.DefaultQuery("fromDate", "")
		toDate := c.DefaultQuery("toDate", "")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":
		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)
		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if fromDate != "" {
			conditions = append(conditions, "report_date >= ?")
			params = append(params, fromDate)
		}

		if toDate != "" {
			conditions = append(conditions, "report_date <= ?")
			params = append(params, toDate)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, "(report_date LIKE ? OR answered LIKE ?)")
			params = append(params, searchTerm, searchTerm)
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}
		countQuery := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM mis_report_tbl
			%s
		`, whereClause)

		var totalCount int
		err := db.QueryRow(countQuery, params...).Scan(&totalCount)
		if err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count records"})
			return
		}

		dataQuery := fmt.Sprintf(`
			SELECT 
			    user_id,
				DATE_FORMAT(report_date, '%%d-%%m-%%Y') as report_date,
				answered,
				not_answered,
				total_calls
			FROM mis_report_tbl
			%s
			ORDER BY report_date DESC
			LIMIT ? OFFSET ?
		`, whereClause)

		queryParams := make([]interface{}, len(params))
		copy(queryParams, params)
		queryParams = append(queryParams, limit, offset)

		rows, err := db.Query(dataQuery, queryParams...)
		if err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
			return
		}
		defer rows.Close()

		var results []models.GetMISReportPayload

		for rows.Next() {
			var item models.GetMISReportPayload
			err := rows.Scan(
				&item.UserId,
				&item.ReportDate,
				&item.Answered,
				&item.NotAnswered,
				&item.TotalsCalls,
			)
			if err != nil {
				utils.LogError(c, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan record"})
				return
			}
			results = append(results, item)
		}

		if err = rows.Err(); err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error iterating rows"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total": totalCount,
			"page":  page,
			"limit": limit,
			"data":  results,
		})
	}
}

func GetTodayReportExcel(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.GetRequestParameters(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":

		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)

		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, `(campaign_title LIKE ? OR sender_mobile LIKE ? OR voice_status LIKE ?)`)
			params = append(params, searchTerm, searchTerm, searchTerm)
		}

		conditions = append(conditions, "send_date = CURDATE()")

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}

		query := fmt.Sprintf(`
			SELECT 
				campaign_title,
				campaign_no,
				sender_mobile,
				temp_category,
				voice_status,
				send_time
			FROM temp_obd_send_tbl_final
			%s
			ORDER BY send_time DESC
		`, whereClause)

		rows, err := db.Query(query, params...)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to fetch records")
			return
		}
		defer rows.Close()

		f := excelize.NewFile()
		sheet := "Today's Report"

		f.SetSheetName("Sheet1", sheet)

		headers := []string{
			"Campaign Title",
			"Campaign No",
			"Sender Mobile",
			"Category",
			"Voice Status",
			"Send Time",
		}

		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheet, cell, header)
		}

		rowIndex := 2

		for rows.Next() {

			var campaignTitle, campaignNo, senderMobile, tempCategory, voiceStatus, sendTime string

			err := rows.Scan(
				&campaignTitle,
				&campaignNo,
				&senderMobile,
				&tempCategory,
				&voiceStatus,
				&sendTime,
			)

			if err != nil {
				utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to scan record")
				return
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), campaignTitle)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), campaignNo)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), senderMobile)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), tempCategory)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), voiceStatus)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), sendTime)

			rowIndex++
		}

		filename := fmt.Sprintf("Today_Report_%s.xlsx", time.Now().Format("2006-01-02_150405"))

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Transfer-Encoding", "binary")

		err = f.Write(c.Writer)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to generate excel")
			return
		}

	}
}

func GetDeliveryReportExcel(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.GetRequestParameters(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		fromDate := c.DefaultQuery("fromDate", "")
		toDate := c.DefaultQuery("toDate", "")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":

		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)

		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if fromDate != "" {
			conditions = append(conditions, "send_date >= ?")
			params = append(params, fromDate)
		}

		if toDate != "" {
			conditions = append(conditions, "send_date <= ?")
			params = append(params, toDate)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, `(campaign_title LIKE ? OR sender_mobile LIKE ? OR voice_status LIKE ?)`)
			params = append(params, searchTerm, searchTerm, searchTerm)
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}

		query := fmt.Sprintf(`
			SELECT 
				campaign_title,
				campaign_no,
				business_no,
				sender_mobile,
				temp_name,
				temp_category,
				voice_status,
				COALESCE(call_duration,' ') AS call_duration,
				send_date,
				send_time,
				count_voice_deduct
			FROM temp_obd_send_tbl_final
			%s
			ORDER BY send_date DESC, send_time DESC
		`, whereClause)

		rows, err := db.Query(query, params...)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to fetch records")
			return
		}
		defer rows.Close()

		f := excelize.NewFile()
		sheet := "Delivery Report"

		f.SetSheetName("Sheet1", sheet)

		headers := []string{
			"Campaign Title",
			"Campaign No",
			"Business No",
			"Sender Mobile",
			"Template Name",
			"Category",
			"Voice Status",
			"Call Duration",
			"Send Date",
			"Send Time",
			"Voice Deduct",
		}

		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheet, cell, header)
		}

		rowIndex := 2

		for rows.Next() {

			var campaignTitle, campaignNo, businessNo, senderMobile string
			var tempName, tempCategory, voiceStatus, callDuration string
			var sendDate, sendTime string
			var countVoiceDeduct int

			err := rows.Scan(
				&campaignTitle,
				&campaignNo,
				&businessNo,
				&senderMobile,
				&tempName,
				&tempCategory,
				&voiceStatus,
				&callDuration,
				&sendDate,
				&sendTime,
				&countVoiceDeduct,
			)

			if err != nil {
				utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to scan record")
				return
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), campaignTitle)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), campaignNo)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), businessNo)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), senderMobile)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), tempName)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), tempCategory)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), voiceStatus)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", rowIndex), callDuration)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", rowIndex), sendDate)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", rowIndex), sendTime)
			f.SetCellValue(sheet, fmt.Sprintf("K%d", rowIndex), countVoiceDeduct)

			rowIndex++
		}

		filename := fmt.Sprintf("Delivery_Report_%s.xlsx", time.Now().Format("2006-01-02_150405"))

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename="+filename)

		err = f.Write(c.Writer)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to generate excel")
			return
		}
	}
}

func GetMISReportExcel(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.GetRequestParameters(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		fromDate := c.DefaultQuery("fromDate", "")
		toDate := c.DefaultQuery("toDate", "")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":

		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)

		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if fromDate != "" {
			conditions = append(conditions, "report_date >= ?")
			params = append(params, fromDate)
		}

		if toDate != "" {
			conditions = append(conditions, "report_date <= ?")
			params = append(params, toDate)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, "(report_date LIKE ? OR answered LIKE ?)")
			params = append(params, searchTerm, searchTerm)
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}

		query := fmt.Sprintf(`
			SELECT 
				report_date,
				total_calls,
				answered,
				not_answered
			FROM mis_report_tbl
			%s
			ORDER BY report_date DESC
		`, whereClause)

		rows, err := db.Query(query, params...)
		if err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
			return
		}
		defer rows.Close()

		f := excelize.NewFile()
		sheet := "MIS Report"

		f.SetSheetName("Sheet1", sheet)

		headers := []string{
			"Report Date",
			"Total Calls",
			"Answered",
			"Not Answered",
		}

		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheet, cell, header)
		}

		rowIndex := 2

		for rows.Next() {

			var reportDate string
			var totalCalls, answered, notAnswered int

			err := rows.Scan(
				&reportDate,
				&totalCalls,
				&answered,
				&notAnswered,
			)

			if err != nil {
				utils.LogError(c, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan record"})
				return
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), reportDate)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), totalCalls)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), answered)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), notAnswered)

			rowIndex++
		}

		filename := fmt.Sprintf("MIS_Report_%s.xlsx", time.Now().Format("2006-01-02_150405"))

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename="+filename)

		err = f.Write(c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate excel"})
			return
		}
	}
}

func GetUsersWiseReportExcel(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.GetRequestParameters(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		fromDate := c.DefaultQuery("fromDate", "")
		toDate := c.DefaultQuery("toDate", "")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":

		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)

		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if fromDate != "" {
			conditions = append(conditions, "report_date >= ?")
			params = append(params, fromDate)
		}

		if toDate != "" {
			conditions = append(conditions, "report_date <= ?")
			params = append(params, toDate)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, "(user_id LIKE ? OR answered LIKE ?)")
			params = append(params, searchTerm, searchTerm)
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}

		query := fmt.Sprintf(`
			SELECT 
				user_id,
				report_date,
				total_calls,
				answered,
				not_answered
			FROM mis_report_tbl
			%s
			ORDER BY report_date DESC
		`, whereClause)

		rows, err := db.Query(query, params...)
		if err != nil {
			utils.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
			return
		}
		defer rows.Close()

		f := excelize.NewFile()
		sheet := "Users Wise Report"

		f.SetSheetName("Sheet1", sheet)

		headers := []string{
			"User ID",
			"Report Date",
			"Total Calls",
			"Answered",
			"Not Answered",
		}

		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheet, cell, header)
		}

		rowIndex := 2

		for rows.Next() {

			var userID string
			var reportDate string
			var totalCalls, answered, notAnswered int

			err := rows.Scan(
				&userID,
				&reportDate,
				&totalCalls,
				&answered,
				&notAnswered,
			)

			if err != nil {
				utils.LogError(c, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan record"})
				return
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), userID)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), reportDate)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), totalCalls)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), answered)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), notAnswered)

			rowIndex++
		}

		filename := fmt.Sprintf("Users_Wise_Report_%s.xlsx", time.Now().Format("2006-01-02_150405"))

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename="+filename)

		err = f.Write(c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate excel"})
			return
		}
	}
}

func GetCampaignWiseReportExcel(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.GetRequestParameters(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {
		case "superadmin":

		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)

		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if search != "" {
			searchTerm := "%" + search + "%"
			conditions = append(conditions, `
				(campaign_title LIKE ? OR 
				 campaign_no LIKE ? OR 
				 business_no LIKE ?)
			`)
			params = append(params, searchTerm, searchTerm, searchTerm)
		}

		var tableName string

		if startDate != "" && endDate != "" {

			tableName = "campaign_wise_report"

			conditions = append(conditions, "send_date BETWEEN ? AND ?")
			params = append(params, startDate, endDate)

		} else {

			tableName = "temp_obd_send_tbl_final"

			conditions = append(conditions, "send_date = CURDATE()")
		}

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		query := fmt.Sprintf(`
			SELECT
				campaign_title,
				campaign_no,
				admin_id,
				user_id,
				business_no,
				temp_category,
				send_date,
				send_time,
				COUNT(CASE WHEN voice_status = 'answered' THEN 1 END) AS answered_calls,
				COUNT(CASE WHEN voice_status = 'failed' THEN 1 END) AS failed_calls,
				COUNT(CASE WHEN voice_status = 'reject' THEN 1 END) AS rejected_calls,
				COUNT(*) AS total_calls,
				voice_id
			FROM %s
			%s
			GROUP BY
				campaign_title,
				campaign_no,
				admin_id,
				user_id,
				business_no,
				temp_category,
				send_date,
				send_time,
				voice_id
			ORDER BY send_time DESC
			LIMIT ? OFFSET ?
		`, tableName, whereClause)

		queryParams := append([]interface{}{}, params...)
		queryParams = append(queryParams, limit, offset)

		rows, err := db.Query(query, queryParams...)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to fetch records")
			return
		}
		defer rows.Close()

		f := excelize.NewFile()
		sheet := "Campaign Report"

		f.SetSheetName("Sheet1", sheet)

		headers := []string{
			"Campaign Title",
			"Campaign No",
			"Admin ID",
			"User ID",
			"Business No",
			"Category",
			"Send Date",
			"Send Time",
			"Answered Calls",
			"Failed Calls",
			"Rejected Calls",
			"Total Calls",
			"Voice ID",
		}

		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheet, cell, header)
		}

		rowIndex := 2

		for rows.Next() {

			var campaignTitle string
			var campaignNo string
			var adminID string
			var userID string
			var businessNo string
			var category string
			var sendDate string
			var sendTime string
			var answered int
			var failed int
			var rejected int
			var total int
			var voiceID int

			err := rows.Scan(
				&campaignTitle,
				&campaignNo,
				&adminID,
				&userID,
				&businessNo,
				&category,
				&sendDate,
				&sendTime,
				&answered,
				&failed,
				&rejected,
				&total,
				&voiceID,
			)

			if err != nil {
				utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to scan record")
				return
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), campaignTitle)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), campaignNo)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), adminID)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), userID)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), businessNo)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), category)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), sendDate)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", rowIndex), sendTime)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", rowIndex), answered)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", rowIndex), failed)
			f.SetCellValue(sheet, fmt.Sprintf("K%d", rowIndex), rejected)
			f.SetCellValue(sheet, fmt.Sprintf("L%d", rowIndex), total)

			rowIndex++
		}

		filename := fmt.Sprintf("Campaign_Report_%s.xlsx", time.Now().Format("2006-01-02_150405"))

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Transfer-Encoding", "binary")

		err = f.Write(c.Writer)
		if err != nil {
			utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to generate excel")
			return
		}

	}
}
func GetCampaignWiseReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.GetRequestParameters(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		search := c.DefaultQuery("search", "")

		var conditions []string
		var params []interface{}

		switch user.Role {

		case "superadmin":

		case "admin":
			conditions = append(conditions, "admin_id = ?")
			params = append(params, user.UserId)

		default:
			conditions = append(conditions, "user_id = ?")
			params = append(params, user.UserId)
		}

		if search != "" {

			searchTerm := "%" + search + "%"

			conditions = append(conditions, `
				(campaign_title LIKE ? OR
				 campaign_no LIKE ? OR
				 business_no LIKE ?)
			`)

			params = append(params, searchTerm, searchTerm, searchTerm)
		}

		var tableName string

		if startDate != "" && endDate != "" {

			tableName = "campaign_wise_report"

			conditions = append(conditions, "send_date BETWEEN ? AND ?")
			params = append(params, startDate, endDate)

		} else {

			tableName = "temp_obd_send_tbl_final"

			conditions = append(conditions, "send_date = CURDATE()")
		}

		whereClause := ""

		if len(conditions) > 0 {
			whereClause = "WHERE " + strings.Join(conditions, " AND ")
		}

		query := fmt.Sprintf(`
			SELECT
				campaign_title,
				campaign_no,
				admin_id,
				user_id,
				business_no,
				temp_category,
				send_date,
				send_time,
				COUNT(CASE WHEN voice_status = 'answered' THEN 1 END) AS answered_calls,
				COUNT(CASE WHEN voice_status = 'failed' THEN 1 END) AS failed_calls,
				COUNT(CASE WHEN voice_status = 'reject' THEN 1 END) AS rejected_calls,
				COUNT(*) AS total_calls,
				voice_id
			FROM %s
			%s
			GROUP BY
				campaign_title,
				campaign_no,
				admin_id,
				user_id,
				business_no,
				temp_category,
				send_date,
				send_time,
				voice_id
			ORDER BY send_time DESC
			LIMIT ? OFFSET ?
		`, tableName, whereClause)

		params = append(params, limit, offset)

		rows, err := db.Query(query, params...)
		if err != nil {
			c.Error(err)
			return
		}
		defer rows.Close()

		type CampaignReportData struct {
			CampaignTitle string `json:"campaignTitle"`
			CampaignNo    string `json:"campaignNo"`
			AdminID       string `json:"adminID"`
			UserID        string `json:"userID"`
			BusinessNo    string `json:"businessNo"`
			Category      string `json:"category"`
			SendDate      string `json:"sendDate"`
			SendTime      string `json:"sendTime"`
			Answered      int    `json:"answered"`
			Failed        int    `json:"failed"`
			Rejected      int    `json:"rejected"`
			Total         int    `json:"total"`
			VoiceID       string `json:"voiceID"`
		}

		campaignReportData := []CampaignReportData{}

		for rows.Next() {
			var row CampaignReportData

			err := rows.Scan(
				&row.CampaignTitle,
				&row.CampaignNo,
				&row.AdminID,
				&row.UserID,
				&row.BusinessNo,
				&row.Category,
				&row.SendDate,
				&row.SendTime,
				&row.Answered,
				&row.Failed,
				&row.Rejected,
				&row.Total,
				&row.VoiceID,
			)

			if err != nil {
				c.Error(err)
				return
			}

			campaignReportData = append(campaignReportData, row)
		}

		// get total count query
		var totalCountData int
		var totalCountQury string = fmt.Sprintf(`SELECT COUNT(*) 
		FROM (
			SELECT
				campaign_title,
				campaign_no,
				admin_id,
				user_id,
				business_no,
				temp_category,
				send_date,
				send_time,
				voice_id
			FROM %s
			%s
			GROUP BY
				campaign_title,
				campaign_no,
				admin_id,
				user_id,
				business_no,
				temp_category,
				send_date,
				send_time,
				voice_id
		) AS total
		`, tableName, whereClause)

		if err := db.QueryRow(totalCountQury, params[:len(params)-2]...).Scan(&totalCountData); err != nil && err != sql.ErrNoRows {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"page":    page,
			"limit":   limit,
			"data":    campaignReportData,
			"total":   totalCountData,
		})
	}
}


// controllers code

import { useEffect, useState } from "react";
import { Search, Download, Filter, ChartColumnIncreasing } from "lucide-react";
import Pagination from "../../Componets/Pagination";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import {GetCamaginWiseReportData }from "../../services/report"
import toast, { Toaster } from 'react-hot-toast';

const TABLE_HEADERS = [
  "Admin Id",
  "User Id",
  "Business NO.",
  "Campaign Title",
  "Campaign No.",
  "Voice Id",
  "Temp Category",
  "Send Date",
  "Recieve Date",
  "Send Time",
  "Answered",
  "Failed",
  "Reject",
  "Total Calls",
];

export default function CampianWiseReport() {
  const [searchTerm, setSearchTerm] = useState("");
  const [fromDate, setFromDate] = useState("");
  const [toDate, setToDate] = useState("");
  const [showDateFilter, setShowDateFilter] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [debouncedSearchTerm, setDebouncedSearchTerm] = useState("");
  const [campaginWiseReportData, setCampaginWiseReportData] = useState([]);
  const [totalDataCount, setTotalDataCount] = useState(0);

  const getCampaginWiseReportData = async () => {
    try{
      let apiResponse = await GetCamaginWiseReportData(fromDate, toDate, currentPage, itemsPerPage, debouncedSearchTerm);
      if (apiResponse.success) {
        console.log(apiResponse, " got my api response")
        setCampaginWiseReportData(apiResponse?.data || [])
        setTotalDataCount(apiResponse?.total || 0)
      }
    }catch(error){
      console.log("Error on api calling: ", error)
    }
  }
  
  useEffect(() => {
    console.log('eheheheheh')
    getCampaginWiseReportData()
  }, [])

  useEffect(()=> {
    const handler = setTimeout(() => {
      setDebouncedSearchTerm(searchTerm);
    }, 500);

    return () => clearTimeout(handler)
  }, [searchTerm])
  
//   const filteredData = campaginWiseReportData.filter((item) => {

//   const matchesSearch =
//     item.mobile.toLowerCase().includes(searchTerm.toLowerCase().trim());

//   const itemDate = new Date(item.sent);

//   const matchesFromDate = fromDate ? itemDate >= fromDate : true;
//   const matchesToDate = toDate ? itemDate <= toDate : true;

//   return matchesSearch && matchesFromDate && matchesToDate;
// });

//   const totalPages = Math.ceil(filteredData.length / itemsPerPage);
//   const paginatedData = filteredData.slice(
//     (currentPage - 1) * itemsPerPage,
//     currentPage * itemsPerPage
//   );

  const handleFilter = async () => {
    try{
      if (fromDate && !toDate) {
        toast.error("Please select to date.")
        return
      }

      if (!fromDate && toDate) {
        toast.error("Please select from date.")
        return
      }

      if (!fromDate && !toDate) {
        toast.error("Please selecte from and to date.")
        return
      }

      await getCampaginWiseReportData();
      setCurrentPage(1);
    }catch(error){
      console.log("Error on applying filteration: ", error)
    }
    
  };

  const handlePageChange = (page) => setCurrentPage(page);
  const handleLimitChange = (value) => { setItemsPerPage(value); setCurrentPage(1); };

  const handleExport = () => {
    const headers = [];
    const csvRows = [
      headers.join(","),
      ...filteredData.map((row, i) =>
        [i + 1, row.sent, row.mobile, row.voiceDuration, row.callDuration, row.noOfPulse, row.revert].join(",")
      ),
    ];
    const blob = new Blob([csvRows.join("\n")], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", `mis-report-${new Date().toISOString().split("T")[0]}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };
  //calender validation
 const today = new Date();
const threeMonthsAgo = new Date();
threeMonthsAgo.setMonth(today.getMonth() - 3);

const yesterday = new Date();
yesterday.setDate(yesterday.getDate() - 1);
   const minusOneDay = (date) => {
  const newDate = new Date(date);
  newDate.setDate(newDate.getDate()-1);
  return newDate;
  
};

  return (
    <div className="flex flex-col min-h-screen bg-slate-50">
      <Toaster />

      {/* Header */}
      <header className="bg-white border-b border-gray-100 px-4 sm:px-6 lg:px-8 py-4 sm:py-5 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 shrink-0 shadow-sm">

        {/* Left */}
        <div className="flex items-center gap-3 min-w-0">
          <div className="w-9 h-9 rounded-lg bg-indigo-50 border border-indigo-100 flex items-center justify-center shrink-0">
            <ChartColumnIncreasing size={18} className="text-indigo-600" />
          </div>
          <div className="min-w-0">
            <h1 className="text-lg font-bold text-gray-900 leading-tight truncate">Campaign Wise Report</h1>
            <p className="text-xs text-gray-400 mt-0.5 whitespace-nowrap">
              {totalDataCount} record{campaginWiseReportData?.length !== 1 ? "s" : ""} found
            </p>
          </div>
        </div>

        {/* Right */}
        <div className="flex flex-wrap items-center gap-2 sm:gap-3">

          {/* Search */}
          <div className="relative flex-grow sm:flex-grow-0 sm:w-56">
            <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none" />
            <input
              type="text"
              placeholder="Search user or mobile..."
              value={searchTerm}
              onChange={(e) => { setSearchTerm(e.target.value); setCurrentPage(1); }}
              className="w-full pl-9 pr-4 py-2 text-sm border border-gray-200 rounded-lg bg-gray-50 hover:bg-white focus:bg-white focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 outline-none transition-all"
            />
          </div>

          {/* Filter Toggle */}
          <button
            onClick={() => setShowDateFilter(!showDateFilter)}
            className={`flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors border whitespace-nowrap
              ${showDateFilter
                ? "bg-indigo-50 border-indigo-200 text-indigo-700"
                : "bg-gray-100 border-gray-200 hover:bg-gray-200 text-gray-700"}`}
          >
            <Filter size={15} />
            Filter
          </button>

          {/* Export */}
          <button
            onClick={handleExport}
            className="flex items-center gap-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition-colors shadow-sm whitespace-nowrap"
          >
            <Download size={15} />
            Export CSV
          </button>
        </div>
      </header>

      {/* Date Filter Bar */}
      {showDateFilter && (
        <div className="px-4 sm:px-6 lg:px-8 py-4 bg-white border-b border-gray-100 relative z-30">
          <div className="flex flex-col sm:flex-row sm:items-end gap-3">

            <div className="flex flex-col gap-1 relative z-40">
              <label className="text-xs font-semibold text-gray-500 uppercase tracking-wide">From Date</label>
              <DatePicker
             selected={fromDate}
             onChange={(date) => setFromDate(date)}
             dateFormat="dd/MM/yyyy"
             placeholderText="Select from date"
             minDate={threeMonthsAgo}
             maxDate={toDate || yesterday}
              popperPlacement="bottom-start"
             popperStrategy="fixed"
             className="w-44 px-3 py-2 text-sm border border-gray-200 rounded-lg bg-gray-50 hover:bg-white focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 outline-none transition-all"
/>
            </div>

            <div className="flex flex-col gap-1 relative z-30">
              <label className="text-xs font-semibold text-gray-500 uppercase tracking-wide">To Date</label>
              <DatePicker
                 selected={toDate}
                 onChange={(date) => setToDate(date)}
                 dateFormat="dd/MM/yyyy"
                 placeholderText="Select to date"
                  maxDate={yesterday}
                  minDate={fromDate || threeMonthsAgo}
                  popperPlacement="bottom-start"
                  popperStrategy="fixed"
                 className="w-44 px-3 py-2 text-sm border border-gray-200 rounded-lg bg-gray-50 hover:bg-white focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 outline-none transition-all"
                 />
            </div>

            <button
              onClick={handleFilter}
              className="px-5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition-colors self-end"
            >
              Apply
            </button>

          </div>
        </div>
      )}

      {/* Table */}
      <div className="flex-1 px-4 sm:px-6 lg:px-8 py-5">
        <div className="bg-white border border-gray-100 rounded-xl shadow-sm overflow-hidden">
          <div className="overflow-x-auto max-h-[420px] overflow-y-auto">
            <table className="w-full text-sm text-left text-gray-700">
              <thead>
                <tr className="bg-indigo-600">
                  {TABLE_HEADERS.map((heading, i) => (
                    <th
                      key={i}
                      className="px-5 py-3 text-[11px] font-semibold uppercase tracking-wider text-white text-center border-r border-indigo-500 last:border-r-0 sticky top-0 bg-indigo-600 z-20"
                    >
                      {heading}
                    </th>
                  ))}
                </tr>
              </thead>

              <tbody className="divide-y divide-gray-50">
                {campaginWiseReportData.length === 0 ? (
                  <tr>
                    <td colSpan={8} className="px-5 py-20 pl-90 sm:pl- items-center text-center">
                      <div className="flex flex-col items-center gap-3 text-gray-400">
                        <div className="w-12 h-12 rounded-full bg-gray-50 flex items-center justify-center">
                          <Search size={20} className="text-gray-300" />
                        </div>
                        <p className="text-sm">No records found</p>
                      </div>
                    </td>
                  </tr>
                ) : (
                  campaginWiseReportData.map((item, idx) => {
                    const isFailed = item.status === "0";
                    return (
                      <tr
                        key={idx}
                        className={`text-center transition-colors hover:bg-indigo-50/40 ${idx % 2 === 0 ? "bg-white" : "bg-gray-50/40"}`}
                      >
                        {/* Admin ID*/}
                        <td className="px-5 py-3.5 border-r border-gray-100">
                          <span className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-indigo-50 text-xs font-medium">
                            {item.adminID}
                          </span>
                        </td>

                        {/* User ID */}
                        <td className="px-5 py-3.5 text-gray-700 font-medium border-r border-gray-100 whitespace-nowrap">
                         {item.userID}
                        </td>

                        {/* Business No. */}
                        <td className="px-5 py-3.5 text-gray-700 border-r border-gray-100">
                          {item.businessNo}
                        </td>

                        {/* Campaign Title */}
                        <td className="px-5 py-3.5 border-r border-gray-100">
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium  ">
                            {item.campaignTitle}
                            <span className="w-1.5 h-1.5  " />
                            
                          </span>
                        </td>

                        {/* Campian No. */}
                        <td className="px-5 py-3.5 border-r border-gray-100">
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium  ">
                          {item.campaignNo}
                          </span>
                        </td>

                        {/* Voice Id */}
                        <td className="px-5 py-3.5 border-r border-gray-100">
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5  text-xs font-medium ">
                            <span className="w-1.5 h-1.5" />
                           {item.voiceID}
                          </span>
                        </td>

                        {/* Temp Category */}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                         {item.category}
                        </td>
                         {/* Send Data */}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                         {item.sendDate}
                        </td>
                         {/* Recieve Data */}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                          {item.sendTime}
                        </td>
                         {/*   Send Time     */}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                         
                        </td>
                         {/*Answered */}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                         
                        </td>
                         {/* Failed */}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                         
                        </td>
                         {/*Reject*/}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                         
                        </td>
                         {/*Total Calls*/}
                        <td className="px-5 py-3.5 font-medium text-gray-800 border-r border-gray-100">
                         
                        </td>
                       </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="border-t border-gray-100">
            <Pagination
              currentPage={currentPage}
              totalPages={Math.ceil(totalDataCount/campaginWiseReportData?.length)}
              totalCount={totalDataCount}
              itemsPerPage={itemsPerPage}
              onPageChange={handlePageChange}
              onLimitChange={handleLimitChange}
            />
          </div>
        </div>
      </div>
    </div>
  );
}