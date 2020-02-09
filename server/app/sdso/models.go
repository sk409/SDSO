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
	scanAllRelation                             = []string{"Project", "User", "Vulnerabilities"}
	teamAllRelation                             = []string{"Projects", "Users"}
	teamUserAllRelation                         = []string{"Team", "User", "Role"}
	teamUserInvitationRequestProjectAllRelation = []string{"TeamUserInvitationRequest", "Project"}
	teamUserInvitationRequestAllRelation        = []string{"InviterUser", "InviteeUser", "Projects", "Role", "Team"}
	testAllRelation                             = []string{"Project", "Results", "Results.Status"}
	testResultAllRelation                       = []string{"Test", "Status"}
)

var (
	meetingRelationUsers = "Users"
)

type branchProtectionRule struct {
	ID         uint   `gorm:"primary_key"`
	Branchname string `gorm:"type:varchar(128);not null"`
	ProjectID  uint   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Project    project
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

func (m meetingMessage) public() interface{} {
	i, err := convert(m)
	if err != nil {
		return m
	}
	ma, ok := i.(map[string]interface{})
	if !ok {
		return m
	}
	u := user{}
	err = first(map[string]interface{}{"id": m.UserID}, &u)
	if err != nil {
		return m
	}
	ma["user"] = u
	if m.ParentID != nil {
		parent := meetingMessage{}
		err = first(map[string]interface{}{"id": *m.ParentID}, &parent)
		if err != nil {
			return m
		}
		ma["parent"] = parent
	}
	return ma
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
	Text string `gorm:"type:text;not null"`
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
	Vulnerabilities []vulnerability `json:"vulnerability"`
}

// func (s scan) public() interface{} {
// 	i, err := convert(s)
// 	if err != nil {
// 		return s
// 	}
// 	m := i.(map[string]interface{})
// 	vulnerabilities := []vulnerability{}
// 	err = find(map[string]interface{}{"scanID": s.ID}, &vulnerabilities)
// 	if err != nil {
// 		m["vulnerabilities"] = []vulnerability{}
// 		return m
// 	}
// 	vp := make([]interface{}, len(vulnerabilities))
// 	for index, v := range vulnerabilities {
// 		p, err := public(v)
// 		if err != nil {
// 			continue
// 		}
// 		vp[index] = p
// 	}
// 	m["vulnerabilities"] = vp
// 	return m
// }

type team struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(256);not null;unique;" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Users     []user    `gorm:"many2many:team_users" json:"users"`
	Projects  []project `json:"projects"`
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

// func (t test) public() interface{} {
// 	i, err := convert(t)
// 	if err != nil {
// 		return t
// 	}
// 	m := i.(map[string]interface{})
// 	testResults := []testResult{}
// 	err = find(map[string]interface{}{"testID": t.ID}, &testResults)
// 	if err != nil {
// 		return t
// 	}
// 	statusText := testStatusSuccessText
// 	if t.Steps != len(testResults) {
// 		statusText = testStatusRunningText
// 	}
// 	for _, testResult := range testResults {
// 		testStatus := testStatus{}
// 		err = first(map[string]interface{}{"id": testResult.TestStatusID}, &testStatus)
// 		if err != nil {
// 			break
// 		}
// 		if testStatus.Text == testStatusFailedText {
// 			statusText = testStatusFailedText
// 		}
// 	}
// 	rp := make([]interface{}, len(testResults))
// 	for index, testResult := range testResults {
// 		if statusText == testStatusSuccessText {
// 			status := testStatus{}
// 			err = first(map[string]interface{}{"id": testResult.TestStatusID}, &status)
// 			if err != nil {
// 				continue
// 			}
// 			statusText = status.Text
// 		}
// 		rp[index] = testResult.public()
// 	}
// 	m["color"] = testStatusColors[statusText]
// 	m["results"] = rp
// 	m["status"] = statusText
// 	return m
// }

type testMessage struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Text      string `gorm:"type:text"`
	TestID    uint   `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	ParentID  *uint
}

// func (t testMessage) public() interface{} {
// 	i, err := convert(t)
// 	if err != nil {
// 		return t
// 	}
// 	m, ok := i.(map[string]interface{})
// 	if !ok {
// 		return t
// 	}
// 	test := test{}
// 	err = first(map[string]interface{}{"id": t.TestID}, &test)
// 	if err != nil {
// 		return t
// 	}
// 	m["test"] = test
// 	u := user{}
// 	err = first(map[string]interface{}{"id": t.UserID}, &u)
// 	if err != nil {
// 		return t
// 	}
// 	m["user"] = u
// 	if t.ParentID != nil {
// 		p := testMessage{}
// 		err = first(map[string]interface{}{"id": *t.ParentID}, &p)
// 		if err != nil {
// 			return t
// 		}
// 		m["parent"] = p
// 	}
// 	return m
// }

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

// func (t testResult) public() interface{} {
// 	i, err := convert(t)
// 	if err != nil {
// 		return t
// 	}
// 	m := i.(map[string]interface{})
// 	ts := testStatus{}
// 	err = first(map[string]interface{}{"id": t.TestStatusID}, &ts)
// 	if err != nil {
// 		return t
// 	}
// 	m["color"] = testStatusColors[ts.Text]
// 	m["status"] = ts.Text
// 	return m
// }

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
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `gorm:"type:varchar(32);not null" json:"name"`
	Description string    `gorm:"type:varchar(128);not null" json:"description"`
	Path        string    `gorm:"type:varchar(256);not null" json:"path"`
	Method      string    `gorm:"type:varchar(8);not null" json:"method"`
	Request     string    `gorm:"type:text;not null" json:"request"`
	Response    string    `gorm:"type:text;not null" json:"response"`
	ProjectID   uint      `gorm:"not null" json:"projectId"`
	ScanID      uint      `gorm:"not null" json:"scanId"`
}
