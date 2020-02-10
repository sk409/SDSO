package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sk409/gocase"
	"github.com/sk409/gotype"
)

var gormDB *gorm.DB

func initGORM() {
	dsn := fmt.Sprintf("root:root@(%s)/sdso?charset=utf8&parseTime=True&loc=Local", databaseHost)
	var err error
	gormDB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	gormDB.AutoMigrate(&user{}).AddUniqueIndex("idx_name_password", "name", "password")
	gormDB.AutoMigrate(&request{})
	gormDB.AutoMigrate(&testStatus{})
	gormDB.AutoMigrate(&team{}).AddForeignKey("founder_user_id", "users(id)", "NO ACTION", "CASCADE")
	gormDB.AutoMigrate(&teamUserRole{})
	gormDB.AutoMigrate(&teamUser{}).AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("role_id", "team_user_roles(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&project{}).AddForeignKey("team_id", "teams(id)", "CASCADE", "CASCADE").AddUniqueIndex("name_team_id_unique", "name", "team_id")
	gormDB.AutoMigrate(&scan{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	gormDB.AutoMigrate(&test{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADe")
	gormDB.AutoMigrate(&testResult{}).AddForeignKey("test_id", "tests(id)", "CASCADE", "CASCADE").AddForeignKey("status_id", "test_statuses(id)", "CASCADE", "CASCADe")
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

func deleteGORM(query interface{}, model interface{}) error {
	return gormDB.Delete(model, query).Error
}

func eagerLoadingGORM(db *gorm.DB, allRelation []string, preloads ...string) *gorm.DB {
	if len(preloads) == 1 && preloads[0] == loadAllRelation {
		preloads = allRelation
	}
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

func findGORM(query interface{}, model interface{}, allRelation []string, preloads ...string) error {
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, allRelation, preloads...)
	err := db.Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func findOrderGORM(query interface{}, order string, model interface{}, allRelation []string, preloads ...string) error {
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, allRelation, preloads...)
	err := db.Order(order).Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func firstGORM(query interface{}, model interface{}, allRelation []string, preloads ...string) error {
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, allRelation, preloads...)
	return db.First(model).Error
}

func saveGORM(query map[string]interface{}, model interface{}) error {
	rv := reflect.ValueOf(model).Elem()
	for key, value := range query {
		fieldname := string(gocase.UpperCamelCase([]byte(key), true))
		fv := rv.FieldByName(fieldname)
		ft := fv.Type()
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.String {
			fv.SetString(value.(string))
		} else if ft.Kind() == reflect.Int {
			if gotype.IsInt(value) {
				fv.SetInt(int64(value.(int)))
			}
		} else if ft.Kind() == reflect.Uint {
			var v uint
			if gotype.IsString(value) {
				s := value.(string)
				ui, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return err
				}
				v = uint(ui)
			} else if gotype.IsUint(value) {
				v = value.(uint)
			} else {
				continue
			}
			if fv.Kind() == reflect.Ptr {
				fv.Set(reflect.ValueOf(&v))
			} else {
				fv.SetUint(uint64(v))
			}
		}
	}
	gormDB.Save(model)
	if gormDB.Error != nil {
		return gormDB.Error
	}
	return nil
}
