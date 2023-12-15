// routers/router.go

package routers

import (
	"bufio"
	"database/sql"
	"fmt"
	"mysqldiff/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		// 设置页面标题为 "SQLDiff"
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "SQLDiff"})
	})

	// 添加接收参数的 API 路由
	router.POST("/api/getDatabases", func(c *gin.Context) {
		var dbInfo models.DatabaseInfo
		if err := c.ShouldBindJSON(&dbInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 在这里处理接收到的 dbInfo 参数，连接数据库并获取所有库名列表
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbInfo.Username, dbInfo.Password, dbInfo.IP, dbInfo.Port))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		rows, err := db.Query("SHOW DATABASES")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var databases []string
		for rows.Next() {
			var database string
			if err := rows.Scan(&database); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			databases = append(databases, database)
		}

		// 返回数据库列表
		c.JSON(http.StatusOK, gin.H{"databases": databases})
	})
	router.POST("/api/uploadSQL", uploadSQLHandler)

	return router
}
func uploadSQLHandler(c *gin.Context) {
	file, err := c.FormFile("sqlFile")
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "No file uploaded"})
		return
	}

	fileName := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, fileName); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Failed to save the uploaded file"})
		return
	}

	// Read SQL file
	sqlFile, err := os.Open(fileName)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Failed to open the uploaded file"})
		return
	}
	defer sqlFile.Close()

	// Extract table names
	tableNames, err := extractTableNames(sqlFile)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Failed to extract table names"})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"tableNames": tableNames})
}

func extractTableNames(file *os.File) ([]string, error) {
	var tableNames []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "CREATE TABLE") {
			tableName := extractTableName(line)
			if tableName != "" {
				tableNames = append(tableNames, tableName)
			}
		}
	}

	return tableNames, nil
}

func extractTableName(statement string) string {
	if strings.Contains(statement, "CREATE TABLE") {
		parts := strings.Fields(statement)
		if len(parts) >= 4 {
			return strings.Trim(parts[3], "(`")
		}
	}
	return ""
}
