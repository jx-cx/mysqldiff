package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// 添加其他字段...
}

func GetUserByID(userID int) (*User, error) {
	// 这里假设你使用了数据库/sql包，需要根据你的实际数据库类型和驱动进行修改
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 执行查询语句
	query := fmt.Sprintf("SELECT id, name FROM users WHERE id = %d", userID)
	row := db.QueryRow(query)

	// 将查询结果映射到 User 结构体
	var user User
	err = row.Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
