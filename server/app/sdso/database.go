package main

import (
	"fmt"
	"reflect"

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
	db.AutoMigrate(&team{})
	db.AutoMigrate(&teamUserRole{})
	db.AutoMigrate(&teamUser{}).AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("role_id", "team_user_roles(id)", "CASCADE", "CASCADE").AddUniqueIndex("team_id_user_id_unique", "team_id", "user_id")
	db.AutoMigrate(&project{}).AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddUniqueIndex("name_team_id_unique", "name", "team_id")
	db.AutoMigrate(&scan{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&test{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADe")
	db.AutoMigrate(&testResult{}).AddForeignKey("test_id", "tests(id)", "CASCADE", "CASCADE").AddForeignKey("test_status_id", "test_statuses(id)", "CASCADE", "CASCADe")
	db.AutoMigrate(&vulnerability{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE").AddForeignKey("scan_id", "scans(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&branchProtectionRule{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&projectUserRole{})
	db.AutoMigrate(&projectUser{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("role_id", "project_user_roles(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&teamUserInvitationRequest{}).AddForeignKey("inviter_user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("invitee_user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddForeignKey("role_id", "team_user_roles(id)", "CASCADE", "CASCADE").AddUniqueIndex("team_id_invitee_user_id_unique", "team_id", "invitee_user_id")
	db.AutoMigrate(&teamUserInvitationRequestProject{}).AddForeignKey("team_user_invitation_request_id", "team_user_invitation_requests(id)", "CASCADE", "CASCADE").AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	insertData()
}

func insertData() {
	insertIfNotExist := func(model interface{}) {
		db.Where(model).First(model)
		rv := reflect.ValueOf(model).Elem()
		id := rv.FieldByName("ID").Uint()
		if id != 0 {
			return
		}
		db.Save(model)
	}
	insertIfNotExist(&testStatus{
		Text: testStatusFailedText,
	})
	insertIfNotExist(&testStatus{
		Text: testStatusRunningText,
	})
	insertIfNotExist(&testStatus{
		Text: testStatusSuccessText,
	})
	insertIfNotExist(&projectUserRole{
		Role: roleProjectUserManager,
	})
	insertIfNotExist(&projectUserRole{
		Role: roleProjectUserMember,
	})
	insertIfNotExist(&teamUserRole{
		Role: roleTeamUserManager,
	})
	insertIfNotExist(&teamUserRole{
		Role: roleTeamUserMember,
	})
}
