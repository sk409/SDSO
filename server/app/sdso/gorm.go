package main

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var gormDB *gorm.DB

func init() {
	dsn := fmt.Sprintf("root:root@(%s)/sdso?charset=utf8&parseTime=True&loc=Local", databaseHost)
	var err error
	gormDB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	gormDB.AutoMigrate(&user{}).AddUniqueIndex("idx_name_password", "name", "password")
	gormDB.AutoMigrate(&request{})
	gormDB.AutoMigrate(&testStatus{})
	gormDB.AutoMigrate(&team{})
	gormDB.AutoMigrate(&teamUserRole{})
	gormDB.AutoMigrate(&teamUser{}).AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("role_id", "team_user_roles(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&project{}).AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddUniqueIndex("name_team_id_unique", "name", "team_id")
	gormDB.AutoMigrate(&scan{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&test{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADe")
	gormDB.AutoMigrate(&testResult{}).AddForeignKey("test_id", "tests(id)", "CASCADE", "CASCADE").AddForeignKey("test_status_id", "test_statuses(id)", "CASCADE", "CASCADe")
	gormDB.AutoMigrate(&vulnerability{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE").AddForeignKey("scan_id", "scans(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&branchProtectionRule{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&projectUserRole{})
	gormDB.AutoMigrate(&projectUser{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("role_id", "project_user_roles(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&teamUserInvitationRequest{}).AddForeignKey("inviter_user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("invitee_user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddForeignKey("role_id", "team_user_roles(id)", "CASCADE", "CASCADE").AddUniqueIndex("team_id_invitee_user_id_unique", "team_id", "invitee_user_id")
	gormDB.AutoMigrate(&teamUserInvitationRequestProject{}).AddForeignKey("team_user_invitation_request_id", "team_user_invitation_requests(id)", "CASCADE", "CASCADE").AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&meeting{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&meetingUser{}).AddForeignKey("meeting_id", "meetings(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddUniqueIndex("meeting_id_user_id_unique", "meeting_id", "user_id")
	gormDB.AutoMigrate(&meetingMessage{}).AddForeignKey("meeting_id", "meetings(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	gormDB.Model(&meetingMessage{}).AddForeignKey("parent_id", "meeting_messages(id)", "SET NULL", "CASCADE")
	gormDB.AutoMigrate(&testMessage{}).AddForeignKey("test_id", "tests(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	gormDB.Model(&testMessage{}).AddForeignKey("parent_id", "test_messages(id)", "SET NULL", "CASCADE")
	gormDB.AutoMigrate(&dastVulnerabilityMessage{}).AddForeignKey("vulnerability_id", "vulnerabilities(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	gormDB.Model(&dastVulnerabilityMessage{}).AddForeignKey("parent_id", "dast_vulnerability_messages(id)", "SET NULL", "CASCADE")
	insertData()
}

func insertData() {
	insertIfNotExist := func(model interface{}) {
		gormDB.Where(model).First(model)
		rv := reflect.ValueOf(model).Elem()
		id := rv.FieldByName("ID").Uint()
		if id != 0 {
			return
		}
		gormDB.Save(model)
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
