package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetCampaignWiseReportExcel(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := utils.GetRequestParameters(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

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

		// Current Date Filter
		conditions = append(conditions, "send_date = CURDATE()")

		whereClause := ""
		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}

		query := fmt.Sprintf(`
			SELECT 
				admin_id,
				user_id
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
		sheet := "Campaign Report"

		f.SetSheetName("Sheet1", sheet)

		headers := []string{
			"Admin ID",
			"User ID",
		}

		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheet, cell, header)
		}

		rowIndex := 2

		for rows.Next() {

			var adminID, userID int

			err := rows.Scan(
				&adminID,
				&userID,
			)

			if err != nil {
				utils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to scan record")
				return
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), adminID)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), userID)

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
