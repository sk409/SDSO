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
	err := findGORM(query, &branchProtectionRules, branchProtectionRuleAllRelation, preloads...)
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
	branchProtectionRule := branchProtectionRule{}
	query := map[string]interface{}{"Branchname": branchname, "ProjectID": projectID}
	err := saveGORM(query, &branchProtectionRule)
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
	err := findGORM(query, &dastVulnerabilityMessages, dastVulnerabilityMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return dastVulnerabilityMessages, nil
}

func (d *dastVulnerabilityMessageRepositoryGORM) findByID(id uint, preloads ...string) (*dastVulnerabilityMessage, error) {
	dastVulnerabilityMessage := dastVulnerabilityMessage{}
	err := firstGORM(map[string]interface{}{"id": id}, &dastVulnerabilityMessage, dastVulnerabilityMessageAllRelation, preloads...)
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
	err := findGORM(query, &meetings, meetingAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func (m *meetingRepositoryGORM) findByID(id uint, preloads ...string) (*meeting, error) {
	meeting := meeting{ID: id}
	err := firstGORM(map[string]interface{}{"id": id}, &meeting, meetingAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (m *meetingRepositoryGORM) findByIDs(ids []uint, preloads ...string) ([]meeting, error) {
	meetings := []meeting{}
	err := findGORM(ids, &meetings, meetingAllRelation, preloads...)
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
	err := findGORM(query, &meetingMessages, meetingMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return meetingMessages, nil
}

func (m *meetingMessageRepositoryGORM) findByID(id uint, preloads ...string) (*meetingMessage, error) {
	meetingMessage := meetingMessage{}
	err := firstGORM(map[string]interface{}{"id": id}, &meetingMessage, meetingMessageAllRelation, preloads...)
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
	err := findGORM(query, &meetingUsers, meetingUserAllRelation, preloads...)
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
	err := findGORM(query, &projects, projectAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (p *projectRepositoryGORM) findByID(id uint, preloads ...string) (*project, error) {
	project := project{}
	err := firstGORM(map[string]interface{}{"id": id}, &project, projectAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *projectRepositoryGORM) findByIDs(ids []uint, preloads ...string) ([]project, error) {
	projects := []project{}
	err := findGORM(ids, &projects, projectAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (p *projectRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*project, error) {
	project := project{}
	err := firstGORM(query, &project, projectAllRelation, preloads...)
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
	err := findGORM(query, &projectUsers, projectUserAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return projectUsers, nil
}

func (p *projectUserRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*projectUser, error) {
	projectUser := projectUser{}
	err := firstGORM(query, &projectUser, projectUserAllRelation, preloads...)
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
	find(map[string]interface{}, ...string) ([]projectUserRole, error)
	findByID(uint, ...string) (*projectUserRole, error)
	findByRole(string, ...string) (*projectUserRole, error)
	save(map[string]interface{}) (*projectUserRole, error)
}

type projectUserRoleRepositoryGORM struct {
}

func (p *projectUserRoleRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]projectUserRole, error) {
	projectUserRoles := []projectUserRole{}
	err := findGORM(query, &projectUserRoles, projectUserRoleAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return projectUserRoles, nil
}

func (p *projectUserRoleRepositoryGORM) findByID(id uint, preloads ...string) (*projectUserRole, error) {
	projectUserRole := projectUserRole{}
	err := firstGORM(map[string]interface{}{"id": id}, &projectUserRole, projectUserRoleAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &projectUserRole, nil
}

func (p *projectUserRoleRepositoryGORM) findByRole(role string, preloads ...string) (*projectUserRole, error) {
	projectUserRole := projectUserRole{}
	err := firstGORM(map[string]interface{}{"role": role}, &projectUserRole, projectUserRoleAllRelation, preloads...)
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
	err := findGORM(query, &scans, scanAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return scans, nil
}

func (s *scanRepositoryGORM) findByID(id uint, preloads ...string) (*scan, error) {
	scan := scan{}
	err := firstGORM(map[string]interface{}{"id": id}, &scan, scanAllRelation, preloads...)
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
	scan := scan{}
	query := map[string]interface{}{"CommitSHA1": commitSHA1, "ProjectID": projectID, "UserID": userID}
	err := saveGORM(query, &scan)
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
	err := findGORM(query, &teams, teamAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (t *teamRepositoryGORM) findByID(id uint, preloads ...string) (*team, error) {
	team := team{ID: id}
	err := firstGORM(map[string]interface{}{"id": id}, &team, teamAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (t *teamRepositoryGORM) findByIDs(ids []uint, preloads ...string) ([]team, error) {
	teams := []team{}
	err := findGORM(ids, &teams, teamAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (t *teamRepositoryGORM) findByName(name string, preloads ...string) (*team, error) {
	team := team{}
	err := firstGORM(map[string]interface{}{"name": name}, &team, teamAllRelation, preloads...)
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
	err := findGORM(query, &teamUsers, teamUserAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return teamUsers, nil
}

func (t *teamUserRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*teamUser, error) {
	teamUser := teamUser{}
	err := firstGORM(query, &teamUser, teamUserAllRelation, preloads...)
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
	return deleteGORM(query, teamUserInvitationRequestProject{})
}

func (t *teamUserInvitationRequestProjectRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]teamUserInvitationRequestProject, error) {
	teamUserInvitationRequestProjects := []teamUserInvitationRequestProject{}
	err := findGORM(query, &teamUserInvitationRequestProjects, teamUserInvitationRequestProjectAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return teamUserInvitationRequestProjects, nil
}

func (t *teamUserInvitationRequestProjectRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*teamUserInvitationRequestProject, error) {
	teamUserInvitationRequestProject := teamUserInvitationRequestProject{}
	err := firstGORM(query, &teamUserInvitationRequestProject, teamUserInvitationRequestProjectAllRelation, preloads...)
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
	return deleteGORM(query, &teamUserInvitationRequest{})
}

func (t *teamUserInvitationRequestRepositoryGORM) deleteByID(id uint) error {
	return deleteGORM(map[string]interface{}{"id": id}, &teamUserInvitationRequest{})
}

func (t *teamUserInvitationRequestRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]teamUserInvitationRequest, error) {
	teamUserInvitationRequests := []teamUserInvitationRequest{}
	err := findGORM(query, &teamUserInvitationRequests, teamUserInvitationRequestAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return teamUserInvitationRequests, nil
}

func (t *teamUserInvitationRequestRepositoryGORM) findByID(id uint, preloads ...string) (*teamUserInvitationRequest, error) {
	teamUserInvitationRequest := teamUserInvitationRequest{}
	err := firstGORM(map[string]interface{}{"id": id}, &teamUserInvitationRequest, teamUserInvitationRequestAllRelation, preloads...)
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
	teamUserInvitationRequest := teamUserInvitationRequest{}
	query := map[string]interface{}{"Message": message, "InviterUserID": inviterUserID, "InviteeUserID": inviteeUserID, "RoleID": roleID, "TeamID": teamID}
	err := saveGORM(query, &teamUserInvitationRequest)
	if err != nil {
		return nil, err
	}
	return &teamUserInvitationRequest, nil
}

type teamUserRoleRepositoryInterface interface {
	findByID(uint, ...string) (*teamUserRole, error)
	findByRole(string, ...string) (*teamUserRole, error)
	saveWith(string) (*teamUserRole, error)
}

type teamUserRoleRepositoryGORM struct {
}

func (t *teamUserRoleRepositoryGORM) findByID(id uint, preloads ...string) (*teamUserRole, error) {
	teamUserRole := teamUserRole{}
	err := firstGORM(map[string]interface{}{"id": id}, &teamUserRole, teamUserRoleAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &teamUserRole, nil
}

func (t *teamUserRoleRepositoryGORM) findByRole(role string, preloads ...string) (*teamUserRole, error) {
	teamUserRole := teamUserRole{}
	err := firstGORM(map[string]interface{}{"role": role}, &teamUserRole, teamUserRoleAllRelation, preloads...)
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
	err := findGORM(query, &tests, testAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func (t *testRepositoryGORM) findByID(id uint, preloads ...string) (*test, error) {
	test := test{}
	err := firstGORM(map[string]interface{}{"id": id}, &test, testAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

func (t *testRepositoryGORM) findOrder(query map[string]interface{}, order string, preloads ...string) ([]test, error) {
	tests := []test{}
	err := findOrderGORM(query, order, &tests, testAllRelation, preloads...)
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
	err := findGORM(query, &testMessages, testMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return testMessages, nil
}

func (t *testMessageRepositoryGORM) findByID(id uint, preloads ...string) (*testMessage, error) {
	testMessage := testMessage{ID: id}
	err := firstGORM(map[string]interface{}{"id": id}, &testMessage, testMessageAllRelation, preloads...)
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
	err := findGORM(query, &testResults, testResultAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return testResults, nil
}

func (t *testResultRepositoryGORM) findByID(id uint, preloads ...string) (*testResult, error) {
	testResult := testResult{}
	err := firstGORM(map[string]interface{}{"id": id}, &testResult, testResultAllRelation, preloads...)
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
	find(map[string]interface{}, ...string) ([]testStatus, error)
	findByID(uint, ...string) (*testStatus, error)
	findByText(string, ...string) (*testStatus, error)
	save(map[string]interface{}) (*testStatus, error)
}

type testStatusRepositoryGORM struct {
}

func (t *testStatusRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]testStatus, error) {
	testStatuses := []testStatus{}
	err := findGORM(query, &testStatuses, testStatusAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return testStatuses, nil
}

func (t *testStatusRepositoryGORM) findByID(id uint, preloads ...string) (*testStatus, error) {
	testStatus := testStatus{ID: id}
	err := firstGORM(map[string]interface{}{"id": id}, &testStatus, testStatusAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &testStatus, nil
}

func (t *testStatusRepositoryGORM) findByText(text string, preloads ...string) (*testStatus, error) {
	testStatus := testStatus{}
	err := firstGORM(map[string]interface{}{"text": text}, &testStatus, testStatusAllRelation, preloads...)
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
	find(map[string]interface{}, ...string) ([]user, error)
	findByID(uint, ...string) (*user, error)
	findByIDs([]uint, ...string) ([]user, error)
	findByName(string, ...string) (*user, error)
	saveWith(name, password, profileImagePath string) (*user, error)
}

type userRepositoryGORM struct {
}

func (u *userRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]user, error) {
	users := []user{}
	err := findGORM(query, &users, userAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepositoryGORM) findByID(id uint, preloads ...string) (*user, error) {
	user := user{}
	err := firstGORM(map[string]interface{}{"id": id}, &user, userAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepositoryGORM) findByIDs(ids []uint, preloads ...string) ([]user, error) {
	users := []user{}
	err := findGORM(ids, &users, userAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepositoryGORM) findByName(name string, preloads ...string) (*user, error) {
	user := user{}
	err := firstGORM(map[string]interface{}{"name": name}, &user, userAllRelation, preloads...)
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
	err := findGORM(query, &vulnerabilities, vulnerabilityAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return vulnerabilities, nil
}

func (v *vulnerabilityRepositoryGORM) findByID(id uint, preloads ...string) (*vulnerability, error) {
	vulnerability := vulnerability{}
	err := firstGORM(map[string]interface{}{"id": id}, &vulnerability, vulnerabilityAllRelation, preloads...)
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
