package main

import (
	"time"
)

var (
	branchProtectionRuleAllRelation             = []string{"Project"}
	dastVulnerabilityMessageAllRelation         = []string{"Parent", "User", "Vulnerability"}
	meetingAllRelation                          = []string{"Project", "Users"}
	meetingMessageAllRelation                   = []string{"Meeting", "Parent", "User"}
	meetingUserAllRelation                      = []string{"Meeting", "User"}
	projectAllRelation                          = []string{"Team", "Users"}
	projectUserAllRelation                      = []string{"Project", "User", "Role"}
	projectUserRoleAllRelation                  = []string{}
	scanAllRelation                             = []string{"Project", "User", "Vulnerabilities", "Vulnerabilities.Scan"}
	teamAllRelation                             = []string{"Projects", "InvitationRequests", "Users"}
	teamUserAllRelation                         = []string{"Team", "User", "Role"}
	teamUserInvitationRequestProjectAllRelation = []string{"TeamUserInvitationRequest", "Project"}
	teamUserInvitationRequestAllRelation        = []string{"InviterUser", "InviteeUser", "Projects", "Role", "Team"}
	teamUserRoleAllRelation                     = []string{}
	testAllRelation                             = []string{"Project", "Results", "Results.Status"}
	testMessageAllRelation                      = []string{"Test", "User", "Parent"}
	testResultAllRelation                       = []string{"Test", "Status"}
	testStatusAllRelation                       = []string{}
	userAllRelation                             = []string{}
	vulnerabilityAllRelation                    = []string{"Project", "Scan"}
)

var (
	meetingRelationProject                     = "Project"
	meetingRelationProjectUsers                = "Project.Users"
	meetingRelationUsers                       = "Users"
	projectRelationTeam                        = "Team"
	projectRelationTeamUsers                   = "Team.Users"
	projectRelationUsers                       = "Users"
	teamRelationInvitationRequests             = "InvitationRequests"
	teamRelationInvitationRequestsInviteeUser  = "InvitationRequests.InviteeUser"
	teamRelationUsers                          = "Users"
	teamUserInvitationRequestRelationTeam      = "Team"
	teamUserInvitationRequestRelationTeamUsers = "Team.Users"
	testRelationProjectUsers                   = "Project.Users"
	testMessageRelationTestProjectUsers        = "Test.Project.Users"
	vulnerabilityRelationProject               = "Project"
	vulnerabilityRelationProjectUsers          = "Project.Users"
)

type branchProtectionRule struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	Branchname string    `gorm:"type:varchar(128);not null" json:"branchname"`
	ProjectID  uint      `gorm:"not null" json:"projectId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Project    project   `json:"project"`
}

type dastVulnerabilityMessage struct {
	ID              uint                      `gorm:"primary_key" json:"id"`
	Text            string                    `gorm:"type:text" json:"text"`
	VulnerabilityID uint                      `gorm:"not null" json:"vulnerabilityId"`
	UserID          uint                      `gorm:"not null" json:"userId"`
	ParentID        *uint                     `json:"parentId"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
	Vulnerability   vulnerability             `json:"vulnerability"`
	User            user                      `json:"user"`
	Parent          *dastVulnerabilityMessage `json:"parent"`
}

type meeting struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(128);not null" json:"name"`
	ProjectID uint      `gorm:"not null" json:"projectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Project   project   `json:"project"`
	Users     []user    `gorm:"many2many:meeting_users" json:"users"`
}

type meetingMessage struct {
	ID        uint            `gorm:"primary_key" json:"id"`
	Text      string          `gorm:"type:text;not null" json:"text"`
	MeetingID uint            `gorm:"not null" json:"meetingId"`
	ParentID  *uint           `json:"parentId"`
	UserID    uint            `gorm:"not null" json:"userId"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Meeting   meeting         `json:"meeting"`
	Parent    *meetingMessage `json:"parent"`
	User      user            `json:"user"`
}

type meetingUser struct {
	MeetingID uint      `gorm:"not null" json:"meetingId"`
	UserID    uint      `gorm:"not null" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Meeting   meeting   `json:"meeting"`
	User      user      `json:"user"`
}

type project struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(128);not null" json:"name"`
	TeamID    uint      `gorm:"not null" json:"teamId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Team      team      `json:"team"`
	Users     []user    `gorm:"many2many:project_users" json:"users"`
}

type projectUser struct {
	ProjectID uint            `gorm:"not null" json:"projectId"`
	UserID    uint            `gorm:"not null" json:"userId"`
	RoleID    uint            `gorm:"not null" json:"roleId"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Project   project         `json:"project"`
	User      user            `json:"user"`
	Role      projectUserRole `json:"role"`
}

type projectUserRole struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Role      string    `gorm:"not null;unique" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type request struct {
	Text string `gorm:"type:text;not null" json:"text"`
}

type scan struct {
	ID              uint            `gorm:"primary_key" json:"id"`
	CommitSHA1      string          `gorm:"type:char(40);not null" json:"commitSha1"`
	ProjectID       uint            `gorm:"not null" json:"projectId"`
	UserID          uint            `gorm:"not null" json:"userId"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	Project         project         `json:"project"`
	User            user            `json:"user"`
	Vulnerabilities []vulnerability `json:"vulnerabilities"`
}

type team struct {
	ID                 uint                        `gorm:"primary_key" json:"id"`
	Name               string                      `gorm:"type:varchar(256);not null;unique;" json:"name"`
	FounderUserID      uint                        `gorm:"not null"`
	CreatedAt          time.Time                   `json:"createdAt"`
	UpdatedAt          time.Time                   `json:"updatedAt"`
	Projects           []project                   `json:"projects"`
	InvitationRequests []teamUserInvitationRequest `json:"invitationRequests"`
	Users              []user                      `gorm:"many2many:team_users" json:"users"`
}

type teamUser struct {
	TeamID    uint         `gorm:"not null" json:"teamId"`
	UserID    uint         `gorm:"not null" json:"userId"`
	RoleID    uint         `gorm:"not null" json:"roleId"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Team      team         `json:"team"`
	User      user         `json:"user"`
	Role      teamUserRole `json:"role"`
}

type teamUserInvitationRequest struct {
	ID            uint         `gorm:"primary_key" json:"id"`
	Message       string       `gorm:"type:varchar(512);" json:"message"`
	InviterUserID uint         `gorm:"not null" json:"inviterUserId"`
	InviteeUserID uint         `gorm:"not null" json:"inviteeUserId"`
	RoleID        uint         `gorm:"not null" json:"roleId"`
	TeamID        uint         `gorm:"not null" json:"teamId"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	InviterUser   user         `json:"inviterUser"`
	InviteeUser   user         `json:"inviteeUser"`
	Projects      []project    `gorm:"many2many:team_user_invitation_request_projects" json:"projects"`
	Role          teamUserRole `json:"role"`
	Team          team         `json:"team"`
}

type teamUserInvitationRequestProject struct {
	TeamUserInvitationRequestID uint                      `gorm:"not null" json:"teamUserInvitationRequestId"`
	ProjectID                   uint                      `gorm:"not null" json:"projectId"`
	CreatedAt                   time.Time                 `json:"createdAt"`
	UpdatedAt                   time.Time                 `json:"updatedAt"`
	TeamUserInvitationRequest   teamUserInvitationRequest `json:"teamUserInvitationRequest"`
	Project                     project                   `json:"project"`
}

type teamUserRole struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Role      string    `gorm:"not null;unique" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type test struct {
	ID         uint         `gorm:"primary_key" json:"id"`
	Steps      int          `gorm:"not null" json:"steps"`
	Branchname string       `gorm:"type:varchar(256); not null;" json:"branchname"`
	CommitSHA1 string       `gorm:"type:char(40);not null;unique" json:"commitSha1"`
	ProjectID  uint         `gorm:"not null" json:"projectId"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
	Project    project      `json:"project"`
	Results    []testResult `json:"results"`
}

type testMessage struct {
	ID        uint         `gorm:"primary_key" json:"id"`
	Text      string       `gorm:"type:text" json:"text"`
	TestID    uint         `gorm:"not null" json:"testId"`
	UserID    uint         `gorm:"not null" json:"userId"`
	ParentID  *uint        `json:"parentId"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Test      test         `json:"test"`
	User      user         `json:"user"`
	Parent    *testMessage `json:"parent"`
}

type testStatus struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Text      string    `gorm:"type:varchar(7);unique" json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type testResult struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Command     string     `gorm:"type:text;not null" json:"command"`
	Output      string     `gorm:"type:text;" json:"output"`
	TestID      uint       `gorm:"not null" json:"testId"`
	StatusID    uint       `gorm:"not null" json:"statusId"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Test        test       `json:"test"`
	Status      testStatus `json:"status"`
}

type user struct {
	ID               uint      `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Name             string    `gorm:"type:varchar(32);not null;unique" json:"name"`
	Password         string    `gorm:"type:char(60);not null;" json:"password"`
	ProfileImagePath string    `gorm:"type:varchar(256);not null;" json:"profileImagePath"`
}

type vulnerability struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(32);not null" json:"name"`
	Description string    `gorm:"type:varchar(128);not null" json:"description"`
	Path        string    `gorm:"type:varchar(256);not null" json:"path"`
	Method      string    `gorm:"type:varchar(8);not null" json:"method"`
	Request     string    `gorm:"type:text;not null" json:"request"`
	Response    string    `gorm:"type:text;not null" json:"response"`
	ProjectID   uint      `gorm:"not null" json:"projectId"`
	ScanID      uint      `gorm:"not null" json:"scanId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Project     project   `json:"project"`
	Scan        scan      `json:"scan"`
}
