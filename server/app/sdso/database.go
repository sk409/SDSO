package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	dsn := fmt.Sprintf("root:root@(%s)/sdso?charset=utf8&parseTime=True&loc=Local", databaseHost)
	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&user{}).AddUniqueIndex("idx_name_password", "name", "password")
	db.AutoMigrate(&request{}, &testStatus{})
	db.AutoMigrate(&project{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&scan{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&test{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADe")
	db.AutoMigrate(&testResult{}).AddForeignKey("test_id", "tests(id)", "CASCADE", "CASCADE").AddForeignKey("test_status_id", "test_statuses(id)", "CASCADE", "CASCADe")
	db.AutoMigrate(&vulnerability{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE").AddForeignKey("scan_id", "scans(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&branchProtectionRule{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	insertData()
}

func insertData() {
	insertIfNotExist := func(text string) {
		count := 0
		db.Model(testStatus{}).Where("text = ?", text).Count(&count)
		if count == 0 {
			testStatus := testStatus{
				Text: text,
			}
			db.Save(&testStatus)
		}
	}
	insertIfNotExist("running")
	insertIfNotExist("success")
	insertIfNotExist("failed")
}
