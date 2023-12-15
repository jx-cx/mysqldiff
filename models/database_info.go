// models/database_info.go

package models

type DatabaseInfo struct {
	IP       string `json:"ip" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
