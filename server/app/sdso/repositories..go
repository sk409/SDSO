package main

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
	teamUserRepository                         teamUserRepositoryInterface
	teamUserInvitationRequestProjectRepository teamUserInvitationRequestProjectRepositoryInterface
	teamUserInvitationRequestRepository        teamUserInvitationRequestRepositoryInterface
	teamUserRoleRepository                     teamUserRoleRepositoryInterface
	testRepository                             testRepositoryInterface
	testMessageRepository                      testMessageRepositoryInterface
	testResultRepository                       testResultRepositoryInterface
	testStatusRepository                       testStatusRepositoryInterface
	userRepository                             userRepositoryInterface
	vulnerabilityRepository                    vulnerabilityRepositoryInterface
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
		teamUserRepository = &teamUserRepositoryGORM{}
		teamUserInvitationRequestRepository = &teamUserInvitationRequestRepositoryGORM{}
		teamUserInvitationRequestProjectRepository = &teamUserInvitationRequestProjectRepositoryGORM{}
		teamUserRoleRepository = &teamUserRoleRepositoryGORM{}
		testRepository = &testRepositoryGORM{}
		testMessageRepository = &testMessageRepositoryGORM{}
		testResultRepository = &testResultRepositoryGORM{}
		testStatusRepository = &testStatusRepositoryGORM{}
		userRepository = &userRepositoryGORM{}
		vulnerabilityRepository = &vulnerabilityRepositoryGORM{}
		initGORM()
	}
}

type branchProtectionRuleRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]branchProtectionRule, error)
	save(map[string]interface{}) (*branchProtectionRule, error)
	saveWith(branchname string, projectID uint) (*branchProtectionRule, error)
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

func (b *branchProtectionRuleRepositoryGORM) save(query map[string]interface{}) (*branchProtectionRule, error) {
	branchProtectionRule := branchProtectionRule{}
	err := saveGORM(query, &branchProtectionRule)
	if err != nil {
		return nil, err
	}
	return &branchProtectionRule, nil
}

func (b *branchProtectionRuleRepositoryGORM) saveWith(branchname string, projectID uint) (*branchProtectionRule, error) {
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
	err := saveGORM(query, &dastVulnerabilityMessage)
	if err != nil {
		return nil, err
	}
	return &dastVulnerabilityMessage, nil
}

type meetingRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meeting, error)
	findByID(uint, ...string) (*meeting, error)
	findByIDs([]uint, ...string) ([]meeting, error)
	first(map[string]interface{}, ...string) (*meeting, error)
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

func (m *meetingRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*meeting, error) {
	meeting := meeting{}
	err := firstGORM(query, &meeting, meetingAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (m *meetingRepositoryGORM) save(query map[string]interface{}) (*meeting, error) {
	meeting := meeting{}
	err := saveGORM(query, &meeting)
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
	err := saveGORM(query, &meetingMessage)
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
	err := saveGORM(query, &meetingUser)
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
	err := saveGORM(query, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

type projectUserRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]projectUser, error)
	first(map[string]interface{}, ...string) (*projectUser, error)
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

func (p *projectUserRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*projectUser, error) {
	projectUser := projectUser{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, projectUserAllRelation, preloads...)
	err := db.First(&projectUser).Error
	if err != nil {
		return nil, err
	}
	return &projectUser, nil
}

func (p *projectUserRepositoryGORM) save(query map[string]interface{}) (*projectUser, error) {
	projectUser := projectUser{}
	err := saveGORM(query, &projectUser)
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
	err := saveGORM(query, &projectUserRole)
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
	err := saveGORM(query, &scan)
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
	first(map[string]interface{}, ...string) (*team, error)
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

func (t *teamRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*team, error) {
	team := team{}
	err := firstGORM(query, &team, teamAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (t *teamRepositoryGORM) save(query map[string]interface{}) (*team, error) {
	team := team{}
	err := saveGORM(query, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

type teamUserRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]teamUser, error)
	first(map[string]interface{}, ...string) (*teamUser, error)
	save(map[string]interface{}) (*teamUser, error)
}

type teamUserRepositoryGORM struct {
}

func (t *teamUserRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]teamUser, error) {
	teamUsers := []teamUser{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, teamUserAllRelation, preloads...)
	err := db.Find(&teamUsers).Error
	if err != nil {
		return nil, err
	}
	return teamUsers, nil
}

func (t *teamUserRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*teamUser, error) {
	teamUser := teamUser{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, teamUserAllRelation, preloads...)
	err := db.First(&teamUser).Error
	if err != nil {
		return nil, err
	}
	return &teamUser, nil
}

func (t *teamUserRepositoryGORM) save(query map[string]interface{}) (*teamUser, error) {
	teamUser := teamUser{}
	err := saveGORM(query, &teamUser)
	if err != nil {
		return nil, err
	}
	return &teamUser, nil
}

type teamUserInvitationRequestProjectRepositoryInterface interface {
	delete(map[string]interface{}) error
	find(map[string]interface{}, ...string) ([]teamUserInvitationRequestProject, error)
	first(map[string]interface{}, ...string) (*teamUserInvitationRequestProject, error)
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

func (t *teamUserInvitationRequestProjectRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*teamUserInvitationRequestProject, error) {
	teamUserInvitationRequestProject := teamUserInvitationRequestProject{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, teamUserInvitationRequestProjectAllRelation, preloads...)
	err := db.First(&teamUserInvitationRequestProject).Error
	if err != nil {
		return nil, err
	}
	return &teamUserInvitationRequestProject, nil
}

func (t *teamUserInvitationRequestProjectRepositoryGORM) save(query map[string]interface{}) (*teamUserInvitationRequestProject, error) {
	teamUserInvitationRequestProject := teamUserInvitationRequestProject{}
	err := saveGORM(query, &teamUserInvitationRequestProject)
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
	first(map[string]interface{}, ...string) (*teamUserInvitationRequest, error)
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

func (t *teamUserInvitationRequestRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*teamUserInvitationRequest, error) {
	teamUserInvitationRequest := teamUserInvitationRequest{}
	err := firstGORM(query, &teamUserInvitationRequest, teamUserInvitationRequestAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &teamUserInvitationRequest, nil
}

func (t *teamUserInvitationRequestRepositoryGORM) save(query map[string]interface{}) (*teamUserInvitationRequest, error) {
	teamUserInvitationRequest := teamUserInvitationRequest{}
	err := saveGORM(query, &teamUserInvitationRequest)
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

type testRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]test, error)
	findByID(uint, ...string) (*test, error)
	findOrder(map[string]interface{}, string, ...string) ([]test, error)
	first(map[string]interface{}, ...string) (*test, error)
	save(map[string]interface{}) (*test, error)
}

type testRepositoryGORM struct {
}

func (t *testRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]test, error) {
	tests := []test{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, testAllRelation, preloads...)
	err := db.Find(&tests).Error
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func (t *testRepositoryGORM) findByID(id uint, preloads ...string) (*test, error) {
	test := test{ID: id}
	db := eagerLoadingGORM(gormDB, testAllRelation, preloads...)
	err := db.First(&test).Error
	if err != nil {
		return nil, err
	}
	return &test, nil
}

func (t *testRepositoryGORM) findOrder(query map[string]interface{}, order string, preloads ...string) ([]test, error) {
	tests := []test{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, testAllRelation, preloads...)
	err := db.Order(order).Find(&tests).Error
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func (t *testRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*test, error) {
	test := test{}
	err := firstGORM(query, &test, testAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

func (t *testRepositoryGORM) save(query map[string]interface{}) (*test, error) {
	test := test{}
	err := saveGORM(query, &test)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

type testMessageRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]testMessage, error)
	findByID(uint, ...string) (*testMessage, error)
	save(map[string]interface{}) (*testMessage, error)
}

type testMessageRepositoryGORM struct {
}

func (t *testMessageRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]testMessage, error) {
	testMessages := []testMessage{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, testMessageAllRelation, preloads...)
	err := db.Find(&testMessages).Error
	if err != nil {
		return nil, err
	}
	return testMessages, nil
}

func (t *testMessageRepositoryGORM) findByID(id uint, preloads ...string) (*testMessage, error) {
	testMessage := testMessage{ID: id}
	db := eagerLoadingGORM(gormDB, testMessageAllRelation, preloads...)
	err := db.First(&testMessage).Error
	if err != nil {
		return nil, err
	}
	return &testMessage, nil
}

func (t *testMessageRepositoryGORM) save(query map[string]interface{}) (*testMessage, error) {
	testMessage := testMessage{}
	err := saveGORM(query, &testMessage)
	if err != nil {
		return nil, err
	}
	return &testMessage, nil
}

type testResultRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]testResult, error)
	findByID(uint, ...string) (*testResult, error)
	save(map[string]interface{}) (*testResult, error)
}

type testResultRepositoryGORM struct {
}

func (t *testResultRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]testResult, error) {
	testResults := []testResult{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, testResultAllRelation, preloads...)
	err := db.Find(&testResults).Error
	if err != nil {
		return nil, err
	}
	return testResults, nil
}

func (t *testResultRepositoryGORM) findByID(id uint, preloads ...string) (*testResult, error) {
	testResult := testResult{ID: id}
	db := eagerLoadingGORM(gormDB, testResultAllRelation, preloads...)
	err := db.First(&testResult).Error
	if err != nil {
		return nil, err
	}
	return &testResult, nil
}

func (t *testResultRepositoryGORM) save(query map[string]interface{}) (*testResult, error) {
	testResult := testResult{}
	err := saveGORM(query, &testResult)
	if err != nil {
		return nil, err
	}
	return &testResult, nil
}

type testStatusRepositoryInterface interface {
	find(map[string]interface{}) ([]testStatus, error)
	findByID(uint) (*testStatus, error)
	findByText(string) (*testStatus, error)
	save(map[string]interface{}) (*testStatus, error)
}

type testStatusRepositoryGORM struct {
}

func (t *testStatusRepositoryGORM) find(query map[string]interface{}) ([]testStatus, error) {
	testStatuses := []testStatus{}
	err := gormDB.Where(query).Find(&testStatuses).Error
	if err != nil {
		return nil, err
	}
	return testStatuses, nil
}

func (t *testStatusRepositoryGORM) findByID(id uint) (*testStatus, error) {
	testStatus := testStatus{ID: id}
	err := gormDB.First(&testStatus).Error
	if err != nil {
		return nil, err
	}
	return &testStatus, nil
}

func (t *testStatusRepositoryGORM) findByText(text string) (*testStatus, error) {
	testStatus := testStatus{}
	err := gormDB.Where("text = ?", text).First(&testStatus).Error
	if err != nil {
		return nil, err
	}
	return &testStatus, nil
}

func (t *testStatusRepositoryGORM) save(query map[string]interface{}) (*testStatus, error) {
	testStatus := testStatus{}
	err := saveGORM(query, &testStatus)
	if err != nil {
		return nil, err
	}
	return &testStatus, nil
}

type userRepositoryInterface interface {
	find(map[string]interface{}) ([]user, error)
	findByID(uint) (*user, error)
	findByIDs([]uint) ([]user, error)
	findByName(string) (*user, error)
	saveWith(name, password, profileImagePath string) (*user, error)
}

type userRepositoryGORM struct {
}

func (u *userRepositoryGORM) find(query map[string]interface{}) ([]user, error) {
	users := []user{}
	err := gormDB.Where(query).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepositoryGORM) findByID(id uint) (*user, error) {
	user := user{ID: id}
	err := gormDB.First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepositoryGORM) findByIDs(ids []uint) ([]user, error) {
	users := []user{}
	err := gormDB.Where(ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
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

type vulnerabilityRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]vulnerability, error)
	findByID(uint, ...string) (*vulnerability, error)
	first(map[string]interface{}, ...string) (*vulnerability, error)
	save(map[string]interface{}) (*vulnerability, error)
}

type vulnerabilityRepositoryGORM struct {
}

func (v *vulnerabilityRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]vulnerability, error) {
	vulnerabilities := []vulnerability{}
	db := gormDB.Where(query)
	db = eagerLoadingGORM(db, vulnerabilityAllRelation, preloads...)
	err := db.Find(&vulnerabilities).Error
	if err != nil {
		return nil, err
	}
	return vulnerabilities, nil
}

func (v *vulnerabilityRepositoryGORM) findByID(id uint, preloads ...string) (*vulnerability, error) {
	vulnerability := vulnerability{ID: id}
	err := findByIDGORM(&vulnerability, vulnerabilityAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &vulnerability, nil
}

func (v *vulnerabilityRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*vulnerability, error) {
	vulnerability := vulnerability{}
	err := firstGORM(query, &vulnerability, vulnerabilityAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &vulnerability, nil
}

func (v *vulnerabilityRepositoryGORM) save(query map[string]interface{}) (*vulnerability, error) {
	vulnerability := vulnerability{}
	err := saveGORM(query, &vulnerability)
	if err != nil {
		return nil, err
	}
	return &vulnerability, nil
}
