package main

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

func (d *dastVulnerabilityMessageRepositoryGORM) findOrderLimit(query map[string]interface{}, order string, limit interface{}, preloads ...string) ([]dastVulnerabilityMessage, error) {
	dastVulnerabilityMessages := []dastVulnerabilityMessage{}
	err := findOrderLimitGORM(query, order, limit, &dastVulnerabilityMessages, dastVulnerabilityMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return dastVulnerabilityMessages, nil
}

func (d *dastVulnerabilityMessageRepositoryGORM) findWhere(where []interface{}, preloads ...string) ([]dastVulnerabilityMessage, error) {
	dastVulnerabilityMessages := []dastVulnerabilityMessage{}
	err := findWhereGORM(where, &dastVulnerabilityMessages, dastVulnerabilityMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return dastVulnerabilityMessages, nil
}

func (d *dastVulnerabilityMessageRepositoryGORM) save(query map[string]interface{}) (*dastVulnerabilityMessage, error) {
	dastVulnerabilityMessage := dastVulnerabilityMessage{}
	err := saveGORM(query, &dastVulnerabilityMessage)
	if err != nil {
		return nil, err
	}
	return &dastVulnerabilityMessage, nil
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

func (m *meetingMessageRepositoryGORM) findOrderLimit(query map[string]interface{}, order string, limit interface{}, preloads ...string) ([]meetingMessage, error) {
	meetingMessages := []meetingMessage{}
	err := findOrderLimitGORM(query, order, limit, &meetingMessages, meetingMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return meetingMessages, nil
}

func (m *meetingMessageRepositoryGORM) findWhere(where []interface{}, preloads ...string) ([]meetingMessage, error) {
	meetingMessages := []meetingMessage{}
	err := findWhereGORM(where, &meetingMessages, meetingMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return meetingMessages, nil
}

func (m *meetingMessageRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*meetingMessage, error) {
	meetingMessage := meetingMessage{}
	err := firstGORM(query, &meetingMessage, meetingMessageAllRelation, preloads...)
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

type meetingMessageViewerRepositoryGORM struct {
}

func (m *meetingMessageViewerRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]meetingMessageViewer, error) {
	meetingMessageViewers := []meetingMessageViewer{}
	err := findGORM(query, &meetingMessageViewers, meetingMessageViewerAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return meetingMessageViewers, nil
}

func (m *meetingMessageViewerRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*meetingMessageViewer, error) {
	meetingMessageViewer := meetingMessageViewer{}
	err := firstGORM(query, &meetingMessageViewer, meetingMessageViewerAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &meetingMessageViewer, nil
}

func (m *meetingMessageViewerRepositoryGORM) save(query map[string]interface{}) (*meetingMessageViewer, error) {
	meetingMessageViewer := meetingMessageViewer{}
	err := saveGORM(query, &meetingMessageViewer)
	if err != nil {
		return nil, err
	}
	return &meetingMessageViewer, nil
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
	teamUserRole := teamUserRole{}
	query := map[string]interface{}{"Role": role}
	err := saveGORM(query, &teamUserRole)
	if err != nil {
		return nil, err
	}
	return &teamUserRole, nil
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

func (t *testMessageRepositoryGORM) findOrderLimit(query map[string]interface{}, order string, limit interface{}, preloads ...string) ([]testMessage, error) {
	testMessages := []testMessage{}
	err := findOrderLimitGORM(query, order, limit, &testMessages, testMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return testMessages, nil
}

func (t *testMessageRepositoryGORM) findWhere(where []interface{}, preloads ...string) ([]testMessage, error) {
	testMessages := []testMessage{}
	err := findWhereGORM(where, &testMessages, testMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return testMessages, nil
}

func (t *testMessageRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*testMessage, error) {
	testMessage := testMessage{}
	err := firstGORM(query, &testMessage, testMessageAllRelation, preloads...)
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

type testMessageViewerRepositoryGORM struct {
}

func (t *testMessageViewerRepositoryGORM) find(query map[string]interface{}, preloads ...string) ([]testMessageViewer, error) {
	testMessageViewers := []testMessageViewer{}
	err := findGORM(query, &testMessageViewers, testMessageAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return testMessageViewers, nil
}

func (t *testMessageViewerRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*testMessageViewer, error) {
	testMessageViewer := testMessageViewer{}
	err := firstGORM(query, &testMessageViewer, testMessageViewerAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &testMessageViewer, nil
}

func (t *testMessageViewerRepositoryGORM) save(query map[string]interface{}) (*testMessageViewer, error) {
	testMessageViewer := testMessageViewer{}
	err := saveGORM(query, &testMessageViewer)
	if err != nil {
		return nil, err
	}
	return &testMessageViewer, nil
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

func (t *testResultRepositoryGORM) update(query map[string]interface{}, values map[string]interface{}) error {
	return updateGORM(query, values, &testResult{})
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

func (u *userRepositoryGORM) first(query map[string]interface{}, preloads ...string) (*user, error) {
	user := user{}
	err := firstGORM(query, &user, userAllRelation, preloads...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepositoryGORM) saveWith(name, password, handlename, email, profileImagePath string) (*user, error) {
	user := user{}
	query := map[string]interface{}{"Name": name, "Password": password, "Handlename": handlename, "Email": email, "ProfileImagePath": profileImagePath}
	err := saveGORM(query, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
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
