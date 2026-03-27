package main

import (
	"blog/config"
	"blog/models"
	"fmt"
)

func main() {
	cfg := config.Load()
	models.InitDB(cfg)

	type ColumnInfo struct {
		Field string
		Type  string
	}

	tables := []string{"sl_article_tags", "sl_directory", "sl_tags"}

	for _, table := range tables {
		fmt.Printf("\n=== %s ===\n", table)
		var cols []ColumnInfo
		models.DB.Raw("SELECT COLUMN_NAME as Field, COLUMN_TYPE as Type FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?", table, "stacklifes").Scan(&cols)
		for _, c := range cols {
			fmt.Printf("%-20s %s\n", c.Field, c.Type)
		}
	}
}
