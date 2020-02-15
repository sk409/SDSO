package main

var (
	branchProtectionRuleRepository             branchProtectionRuleRepositoryInterface
	dastVulnerabilityMessageRepository         dastVulnerabilityMessageRepositoryInterface
	meetingRepository                          meetingRepositoryInterface
	meetingMessageRepository                   meetingMessageRepositoryInterface
	meetingMessageViewerRepository             meetingMessageViewerRepositoryInterface
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
	testMessageViewerRepository                testMessageViewerRepositoryInterface
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
		meetingMessageViewerRepository = &meetingMessageViewerRepositoryGORM{}
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
		testMessageViewerRepository = &testMessageViewerRepositoryGORM{}
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

type dastVulnerabilityMessageRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]dastVulnerabilityMessage, error)
	findByID(uint, ...string) (*dastVulnerabilityMessage, error)
	findOrderLimit(query map[string]interface{}, order string, limit interface{}, preloads ...string) ([]dastVulnerabilityMessage, error)
	findWhere([]interface{}, ...string) ([]dastVulnerabilityMessage, error)
	save(map[string]interface{}) (*dastVulnerabilityMessage, error)
}

type meetingRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meeting, error)
	findByID(uint, ...string) (*meeting, error)
	findByIDs([]uint, ...string) ([]meeting, error)
	first(map[string]interface{}, ...string) (*meeting, error)
	save(map[string]interface{}) (*meeting, error)
}

type meetingMessageRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meetingMessage, error)
	findByID(uint, ...string) (*meetingMessage, error)
	findOrderLimit(query map[string]interface{}, order string, limit interface{}, preloads ...string) ([]meetingMessage, error)
	findWhere([]interface{}, ...string) ([]meetingMessage, error)
	first(map[string]interface{}, ...string) (*meetingMessage, error)
	save(map[string]interface{}) (*meetingMessage, error)
}

type meetingMessageViewerRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meetingMessageViewer, error)
	first(map[string]interface{}, ...string) (*meetingMessageViewer, error)
	save(map[string]interface{}) (*meetingMessageViewer, error)
}

type meetingUserRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]meetingUser, error)
	save(map[string]interface{}) (*meetingUser, error)
}

type projectRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]project, error)
	findByID(uint, ...string) (*project, error)
	findByIDs([]uint, ...string) ([]project, error)
	first(map[string]interface{}, ...string) (*project, error)
	save(map[string]interface{}) (*project, error)
}

type projectUserRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]projectUser, error)
	first(map[string]interface{}, ...string) (*projectUser, error)
	save(map[string]interface{}) (*projectUser, error)
}

type projectUserRoleRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]projectUserRole, error)
	findByID(uint, ...string) (*projectUserRole, error)
	findByRole(string, ...string) (*projectUserRole, error)
	save(map[string]interface{}) (*projectUserRole, error)
}

type sacnRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]scan, error)
	findByID(uint, ...string) (*scan, error)
	save(map[string]interface{}) (*scan, error)
	saveWith(commitSHA1 string, projectID, userID uint) (*scan, error)
}

type teamRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]team, error)
	findByID(uint, ...string) (*team, error)
	findByIDs([]uint, ...string) ([]team, error)
	findByName(string, ...string) (*team, error)
	first(map[string]interface{}, ...string) (*team, error)
	save(map[string]interface{}) (*team, error)
}

type teamUserRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]teamUser, error)
	first(map[string]interface{}, ...string) (*teamUser, error)
	save(map[string]interface{}) (*teamUser, error)
}

type teamUserInvitationRequestProjectRepositoryInterface interface {
	delete(map[string]interface{}) error
	find(map[string]interface{}, ...string) ([]teamUserInvitationRequestProject, error)
	first(map[string]interface{}, ...string) (*teamUserInvitationRequestProject, error)
	save(map[string]interface{}) (*teamUserInvitationRequestProject, error)
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

type teamUserRoleRepositoryInterface interface {
	findByID(uint, ...string) (*teamUserRole, error)
	findByRole(string, ...string) (*teamUserRole, error)
	saveWith(string) (*teamUserRole, error)
}

type testRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]test, error)
	findByID(uint, ...string) (*test, error)
	findOrder(map[string]interface{}, string, ...string) ([]test, error)
	first(map[string]interface{}, ...string) (*test, error)
	save(map[string]interface{}) (*test, error)
}

type testMessageRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]testMessage, error)
	findByID(uint, ...string) (*testMessage, error)
	findOrderLimit(query map[string]interface{}, order string, limit interface{}, preloads ...string) ([]testMessage, error)
	findWhere([]interface{}, ...string) ([]testMessage, error)
	first(map[string]interface{}, ...string) (*testMessage, error)
	save(map[string]interface{}) (*testMessage, error)
}

type testMessageViewerRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]testMessageViewer, error)
	first(map[string]interface{}, ...string) (*testMessageViewer, error)
	save(map[string]interface{}) (*testMessageViewer, error)
}

type testResultRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]testResult, error)
	findByID(uint, ...string) (*testResult, error)
	save(map[string]interface{}) (*testResult, error)
	update(map[string]interface{}, map[string]interface{}) error
}

type testStatusRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]testStatus, error)
	findByID(uint, ...string) (*testStatus, error)
	findByText(string, ...string) (*testStatus, error)
	save(map[string]interface{}) (*testStatus, error)
}

type userRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]user, error)
	findByID(uint, ...string) (*user, error)
	findByIDs([]uint, ...string) ([]user, error)
	findByName(string, ...string) (*user, error)
	first(map[string]interface{}, ...string) (*user, error)
	saveWith(name, password, handlename, email, profileImagePath string) (*user, error)
}

type vulnerabilityRepositoryInterface interface {
	find(map[string]interface{}, ...string) ([]vulnerability, error)
	findByID(uint, ...string) (*vulnerability, error)
	first(map[string]interface{}, ...string) (*vulnerability, error)
	save(map[string]interface{}) (*vulnerability, error)
}
