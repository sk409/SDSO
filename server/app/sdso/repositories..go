package main

import (
	"github.com/jinzhu/gorm"
)

var (
	branchProtectionRuleRepository             branchProtectionRuleRepositoryInterface
	dastVulnerabilityMessageRepository         dastVulnerabilityMessageRepositoryInterface
	meetingRepository                          meetingRepositoryInterface
	meetingMessageRepository                   meetingMessageRepositoryInterface
	meetingUserRepository                      meetingUserRepositoryInterface
	projectRepository                          projectRepositoryInterface
	projectUserRepository                      projectUserRepositoryInterface
	projectUserRoleRepository                  projectUserRoleRepositoryInterface
	scanRepository                             sacnRepositoryInterface
	teamRepository                             teamRepositoryInterface
	teamUserInvitationRequestProjectRepository teamUserInvitationRequestProjectRepositoryInterface
	teamUserInvitationRequestRepository        teamUserInvitationRequestRepositoryInterface
	teamUserRoleRepository                     teamUserRoleRepositoryInterface
	userRepository                             userRepositoryInterface
)

const (
	storageTypeGORM = "storageTypeGORM"
	storageType     = storageTypeGORM
	loadAllRelation = "__loadAllRelation"
)

func init() {
	switch storageType {
	case storageTypeGORM:
		branchProtectionRuleRepository = &branchProtectionRuleRepositoryGORM{}
		dastVulnerabilityMessageRepository = &dastVulnerabilityMessageRepositoryGORM{}
		meetingRepository = &meetingRepositoryGORM{}
		meetingMessageRepository = &meetingMessageRepositoryGORM{}
		meetingUserRepository = &meetingUserRepositoryGORM{}
		projectRepository = &projectRepositoryGORM{}
		projectUserRepository = &projectUserRepositoryGORM{}
		projectUserRoleRepository = &projectUserRoleRepositoryGORM{}
		scanRepository = &scanRepositoryGORM{}
		teamRepository = &teamRepositoryGORM{}
		teamUserInvitationRequestRepository = &teamUserInvitationRequestRepositoryGORM{}
		teamUserInvitationRequestProjectRepository = &teamUserInvitationRequestProjectRepositoryGORM{}
		teamUserRoleRepository = &teamUserRoleRepositoryGORM{}
		userRepository = &userRepositoryGORM{}
	}
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

type branchProtectionRuleRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]branchProtectionRule, error)
	save(string, uint) (*branchProtectionRule, error)
}

type branchProtectionRuleRepositoryGORM struct {
}

func (b *branchProtectionRuleRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]branchProtectionRule, error) {
	branchProtectionRules := []branchProtectionRule{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, branchProtectionRuleAllRelation, preloads...)
	err := db.Find(&branchProtectionRules).Error
	if err != nil {
		return nil, err
	}
	return branchProtectionRules, nil
}

func (b *branchProtectionRuleRepositoryGORM) save(branchname string, projectID uint) (*branchProtectionRule, error) {
	branchProtectionRule := branchProtectionRule{Branchname: branchname, ProjectID: projectID}
	err := gormDB.Save(&branchProtectionRule).Error
	if err != nil {
		return nil, err
	}
	return &branchProtectionRule, nil
}

type dastVulnerabilityMessageRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]dastVulnerabilityMessage, error)
	findByID(uint, ...string) (*dastVulnerabilityMessage, error)
	save(map[string]interface{}) (*dastVulnerabilityMessage, error)
}

type dastVulnerabilityMessageRepositoryGORM struct {
}

func (d *dastVulnerabilityMessageRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]dastVulnerabilityMessage, error) {
	dastVulnerabilityMessages := []dastVulnerabilityMessage{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, dastVulnerabilityMessageAllRelation, preloads...)
	err := db.Find(&dastVulnerabilityMessages).Error
	if err != nil {
		return nil, err
	}
	return dastVulnerabilityMessages, nil
}

func (d *dastVulnerabilityMessageRepositoryGORM) findByID(id uint, preloads ...string) (*dastVulnerabilityMessage, error) {
	dastVulnerabilityMessage := dastVulnerabilityMessage{ID: id}
	db := gormDB
	db = eagerLoadingGORM(db, dastVulnerabilityMessageAllRelation, preloads...)
	err := db.First(&dastVulnerabilityMessage).Error
	if err != nil {
		return nil, err
	}
	return &dastVulnerabilityMessage, nil
}

func (d *dastVulnerabilityMessageRepositoryGORM) save(query map[string]interface{}) (*dastVulnerabilityMessage, error) {
	dastVulnerabilityMessage := dastVulnerabilityMessage{}
	err := save(query, &dastVulnerabilityMessage)
	if err != nil {
		return nil, err
	}
	return &dastVulnerabilityMessage, nil
}

type meetingRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meeting, error)
	findByID(uint, ...string) (*meeting, error)
	findByIDs([]uint, ...string) ([]meeting, error)
	save(map[string]interface{}) (*meeting, error)
}

type meetingRepositoryGORM struct {
}

func (m *meetingRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]meeting, error) {
	meetings := []meeting{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, meetingAllRelation, preloads...)
	err := db.Find(&meetings).Error
	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func (m *meetingRepositoryGORM) findByID(id uint, preloads ...string) (*meeting, error) {
	meeting := meeting{ID: id}
	db := eagerLoadingGORM(gormDB, meetingAllRelation, preloads...)
	err := db.First(&meeting).Error
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (m *meetingRepositoryGORM) findByIDs(ids []uint, preloads ...string) ([]meeting, error) {
	meetings := []meeting{}
	db := gormDB.Where(ids)
	db = eagerLoadingGORM(db, meetingAllRelation, preloads...)
	err := db.Find(&meetings).Error
	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func (m *meetingRepositoryGORM) save(query map[string]interface{}) (*meeting, error) {
	meeting := meeting{}
	err := save(query, &meeting)
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

type meetingMessageRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meetingMessage, error)
	findByID(uint, ...string) (*meetingMessage, error)
	save(map[string]interface{}) (*meetingMessage, error)
}

type meetingMessageRepositoryGORM struct {
}

func (m *meetingMessageRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]meetingMessage, error) {
	meetingMessages := []meetingMessage{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, meetingMessageAllRelation, preloads...)
	err := db.Find(&meetingMessages).Error
	if err != nil {
		return nil, err
	}
	return meetingMessages, nil
}

func (m *meetingMessageRepositoryGORM) findByID(id uint, preloads ...string) (*meetingMessage, error) {
	meetingMessage := meetingMessage{ID: id}
	db := eagerLoadingGORM(gormDB, meetingMessageAllRelation, preloads...)
	err := db.First(&meetingMessage).Error
	if err != nil {
		return nil, err
	}
	return &meetingMessage, nil
}

func (m *meetingMessageRepositoryGORM) save(query map[string]interface{}) (*meetingMessage, error) {
	meetingMessage := meetingMessage{}
	err := save(query, &meetingMessage)
	if err != nil {
		return nil, err
	}
	return &meetingMessage, nil
}

type meetingUserRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meetingUser, error)
	save(map[string]interface{}) (*meetingUser, error)
}

type meetingUserRepositoryGORM struct {
}

func (m *meetingUserRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]meetingUser, error) {
	meetingUsers := []meetingUser{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, meetingUserAllRelation, preloads...)
	err := db.Find(&meetingUsers).Error
	if err != nil {
		return nil, err
	}
	return meetingUsers, nil
}

func (m *meetingUserRepositoryGORM) save(query map[string]interface{}) (*meetingUser, error) {
	meetingUser := meetingUser{}
	err := save(query, &meetingUser)
	if err != nil {
		return nil, err
	}
	return &meetingUser, nil
}

type projectRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]project, error)
	findByID(uint, ...string) (*project, error)
	findByIDs([]uint, ...string) ([]project, error)
	first(map[string]interface{}, ...string) (*project, error)
	save(map[string]interface{}) (*project, error)
}

type projectRepositoryGORM struct {
}

func (p *projectRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]project, error) {
	projects := []project{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, projectAllRelation, preloads...)
	err := db.Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (p *projectRepositoryGORM) findByID(id uint, preloads ...string) (*project, error) {
	project := project{ID: id}
	db := eagerLoadingGORM(gormDB, projectAllRelation, preloads...)
	err := db.First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *projectRepositoryGORM) findByIDs(ids []uint, preloads ...string) ([]project, error) {
	projects := []project{}
	db := gormDB.Where(ids)
	db = eagerLoadingGORM(db, projectAllRelation, preloads...)
	err := db.First(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (p *projectRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*project, error) {
	project := project{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, projectAllRelation, preloads...)
	err := db.First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *projectRepositoryGORM) save(query map[string]interface{}) (*project, error) {
	project := project{}
	err := save(query, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

type projectUserRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]projectUser, error)
	save(map[string]interface{}) (*projectUser, error)
}

type projectUserRepositoryGORM struct {
}

func (p *projectUserRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]projectUser, error) {
	projectUsers := []projectUser{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, projectUserAllRelation, preloads...)
	err := db.Find(&projectUsers).Error
	if err != nil {
		return nil, err
	}
	return projectUsers, nil
}

func (p *projectUserRepositoryGORM) save(query map[string]interface{}) (*projectUser, error) {
	projectUser := projectUser{}
	err := save(query, &projectUser)
	if err != nil {
		return nil, err
	}
	return &projectUser, nil
}

type projectUserRoleRepositoryInterface interface {
	find(map[string]interface{}) ([]projectUserRole, error)
	findByID(uint) (*projectUserRole, error)
	findByRole(string) (*projectUserRole, error)
	save(map[string]interface{}) (*projectUserRole, error)
}

type projectUserRoleRepositoryGORM struct {
}

func (p *projectUserRoleRepositoryGORM) find(query map[string]interface{}) ([]projectUserRole, error) {
	projectUserRoles := []projectUserRole{}
	err := gormDB.Where(query).Find(&projectUserRoles).Error
	if err != nil {
		return nil, err
	}
	return projectUserRoles, nil
}

func (p *projectUserRoleRepositoryGORM) findByID(id uint) (*projectUserRole, error) {
	projectUserRole := projectUserRole{ID: id}
	err := gormDB.First(&projectUserRole).Error
	if err != nil {
		return nil, err
	}
	return &projectUserRole, nil
}

func (p *projectUserRoleRepositoryGORM) findByRole(role string) (*projectUserRole, error) {
	projectUserRole := projectUserRole{}
	err := gormDB.Where("role = ?", role).First(&projectUserRole).Error
	if err != nil {
		return nil, err
	}
	return &projectUserRole, nil
}

func (p *projectUserRoleRepositoryGORM) save(query map[string]interface{}) (*projectUserRole, error) {
	projectUserRole := projectUserRole{}
	err := save(query, &projectUserRole)
	if err != nil {
		return nil, err
	}
	return &projectUserRole, nil
}

type sacnRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]scan, error)
	findByID(uint, ...string) (*scan, error)
	save(map[string]interface{}) (*scan, error)
	saveWith(commitSHA1 string, projectID, userID uint) (*scan, error)
}

type scanRepositoryGORM struct {
}

func (s *scanRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]scan, error) {
	scans := []scan{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, scanAllRelation, preloads...)
	err := db.Find(&scans).Error
	if err != nil {
		return nil, err
	}
	return scans, nil
}

func (s *scanRepositoryGORM) findByID(id uint, preloads ...string) (*scan, error) {
	scan := scan{ID: id}
	db := eagerLoadingGORM(gormDB, scanAllRelation, preloads...)
	err := db.First(&scan).Error
	if err != nil {
		return nil, err
	}
	return &scan, nil
}

func (s *scanRepositoryGORM) save(query map[string]interface{}) (*scan, error) {
	scan := scan{}
	err := save(query, &scan)
	if err != nil {
		return nil, err
	}
	return &scan, nil
}

func (s *scanRepositoryGORM) saveWith(commitSHA1 string, projectID, userID uint) (*scan, error) {
	scan := scan{CommitSHA1: commitSHA1, ProjectID: projectID, UserID: userID}
	err := gormDB.Save(&scan).Error
	if err != nil {
		return nil, err
	}
	return &scan, nil
}

type teamRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]team, error)
	findByID(uint, ...string) (*team, error)
	findByIDs([]uint, ...string) ([]team, error)
	findByName(string, ...string) (*team, error)
	save(map[string]interface{}) (*team, error)
}

type teamRepositoryGORM struct {
}

func (t *teamRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]team, error) {
	teams := []team{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, teamAllRelation, preloads...)
	err := db.Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (t *teamRepositoryGORM) findByID(id uint, preloads ...string) (*team, error) {
	team := team{ID: id}
	db := eagerLoadingGORM(gormDB, teamAllRelation, preloads...)
	err := db.First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (t *teamRepositoryGORM) findByIDs(ids []uint, preloads ...string) ([]team, error) {
	teams := []team{}
	db := gormDB.Where(ids)
	db = eagerLoadingGORM(db, teamAllRelation, preloads...)
	err := db.Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (t *teamRepositoryGORM) findByName(name string, preloads ...string) (*team, error) {
	team := team{}
	db := gormDB.Where("name = ?", name)
	db = eagerLoadingGORM(db, teamAllRelation, preloads...)
	err := db.First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (t *teamRepositoryGORM) save(query map[string]interface{}) (*team, error) {
	team := team{}
	err := save(query, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

type teamUserInvitationRequestProjectRepositoryInterface interface {
	delete(map[string]interface{}) error
	find(map[string]interface{}, ...string) ([]teamUserInvitationRequestProject, error)
	save(map[string]interface{}) (*teamUserInvitationRequestProject, error)
}

type teamUserInvitationRequestProjectRepositoryGORM struct {
}

func (t *teamUserInvitationRequestProjectRepositoryGORM) delete(query map[string]interface{}) error {
	return gormDB.Delete(&teamUserInvitationRequestProject{}, query).Error
}

func (t *teamUserInvitationRequestProjectRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]teamUserInvitationRequestProject, error) {
	teamUserInvitationRequestProjects := []teamUserInvitationRequestProject{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, teamUserInvitationRequestProjectAllRelation, preloads...)
	err := db.Find(&teamUserInvitationRequestProjects).Error
	if err != nil {
		return nil, err
	}
	return teamUserInvitationRequestProjects, nil
}

func (t *teamUserInvitationRequestProjectRepositoryGORM) save(query map[string]interface{}) (*teamUserInvitationRequestProject, error) {
	teamUserInvitationRequestProject := teamUserInvitationRequestProject{}
	err := save(query, &teamUserInvitationRequestProject)
	if err != nil {
		return nil, err
	}
	return &teamUserInvitationRequestProject, nil
}

type teamUserInvitationRequestRepositoryInterface interface {
	delete(map[string]interface{}) error
	deleteByID(uint) error
	find(map[string]interface{}, ...string) ([]teamUserInvitationRequest, error)
	findByID(uint, ...string) (*teamUserInvitationRequest, error)
	save(map[string]interface{}) (*teamUserInvitationRequest, error)
	saveWith(string, uint, uint, uint, uint) (*teamUserInvitationRequest, error)
}

type teamUserInvitationRequestRepositoryGORM struct {
}

func (t *teamUserInvitationRequestRepositoryGORM) delete(query map[string]interface{}) error {
	return gormDB.Delete(&teamUserInvitationRequest{}, query).Error
}

func (t *teamUserInvitationRequestRepositoryGORM) deleteByID(id uint) error {
	return gormDB.Delete(&teamUserInvitationRequest{}, "id = ?", id).Error
}

func (t *teamUserInvitationRequestRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]teamUserInvitationRequest, error) {
	teamUserInvitationRequests := []teamUserInvitationRequest{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, teamUserInvitationRequestAllRelation, preloads...)
	err := db.Find(&teamUserInvitationRequests).Error
	if err != nil {
		return nil, err
	}
	return teamUserInvitationRequests, nil
}

func (t *teamUserInvitationRequestRepositoryGORM) findByID(id uint, preloads ...string) (*teamUserInvitationRequest, error) {
	teamUserInvitationRequest := teamUserInvitationRequest{ID: id}
	db := eagerLoadingGORM(gormDB, teamUserInvitationRequestAllRelation, preloads...)
	err := db.First(&teamUserInvitationRequest).Error
	if err != nil {
		return nil, err
	}
	return &teamUserInvitationRequest, nil
}

func (t *teamUserInvitationRequestRepositoryGORM) save(query map[string]interface{}) (*teamUserInvitationRequest, error) {
	teamUserInvitationRequest := teamUserInvitationRequest{}
	err := save(query, &teamUserInvitationRequest)
	if err != nil {
		return nil, err
	}
	return &teamUserInvitationRequest, nil
}

func (t *teamUserInvitationRequestRepositoryGORM) saveWith(message string, inviterUserID uint, inviteeUserID uint, roleID uint, teamID uint) (*teamUserInvitationRequest, error) {
	teamUserInvitationRequest := teamUserInvitationRequest{Message: message, InviterUserID: inviterUserID, InviteeUserID: inviteeUserID, RoleID: roleID, TeamID: teamID}
	err := gormDB.Save(&teamUserInvitationRequest).Error
	if err != nil {
		return nil, err
	}
	return &teamUserInvitationRequest, nil
}

type teamUserRoleRepositoryInterface interface {
	findByID(uint) (*teamUserRole, error)
	findByRole(string) (*teamUserRole, error)
	saveWith(string) (*teamUserRole, error)
}

type teamUserRoleRepositoryGORM struct {
}

func (t *teamUserRoleRepositoryGORM) findByID(id uint) (*teamUserRole, error) {
	teamUserRole := teamUserRole{ID: id}
	err := gormDB.First(&teamUserRole).Error
	if err != nil {
		return nil, err
	}
	return &teamUserRole, nil
}

func (t *teamUserRoleRepositoryGORM) findByRole(role string) (*teamUserRole, error) {
	teamUserRole := teamUserRole{}
	err := gormDB.Where("role = ?", role).First(&teamUserRole).Error
	if err != nil {
		return nil, err
	}
	return &teamUserRole, nil
}

func (t *teamUserRoleRepositoryGORM) saveWith(role string) (*teamUserRole, error) {
	teamUserRole := teamUserRole{Role: role}
	err := gormDB.Save(&teamUserRole).Error
	if err != nil {
		return nil, err
	}
	return &teamUserRole, nil
}

type userRepositoryInterface interface {
	findByName(name string) (*user, error)
	saveWith(name, password, profileImagePath string) (*user, error)
}

type userRepositoryGORM struct {
}

func (u *userRepositoryGORM) findByName(name string) (*user, error) {
	user := user{}
	err := gormDB.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepositoryGORM) saveWith(name, password, profileImagePath string) (*user, error) {
	user := user{Name: name, Password: password, ProfileImagePath: profileImagePath}
	err := gormDB.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
