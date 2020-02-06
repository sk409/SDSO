package main

import (
	"time"
)

type branchProtectionRule struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	BranchName string `gorm:"type:varchar(128);not null"`
	ProjectID  uint   `gorm:"not null"`
}

type build struct {
	Docker []docker
	Steps  []interface{}
}

type commit struct {
	Branchname string
	Date       string
	Diff       string
	Message    string
	SHA1       string
}

type config struct {
	Version int
	Jobs    jobs
}

type docker struct {
	Image string
}

type jobs struct {
	Build build
}

type meeting struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(128);not null"`
	ProjectID uint   `gorm:"not null"`
}

func (m meeting) public() interface{} {
	i, err := convert(m)
	if err != nil {
		return m
	}
	ma, ok := i.(map[string]interface{})
	if !ok {
		return m
	}
	meetingUsers := []meetingUser{}
	err = find(map[string]interface{}{"meetingID": m.ID}, &meetingUsers)
	if err != nil {
		return m
	}
	userIDs := make([]uint, len(meetingUsers))
	for index, meetingUser := range meetingUsers {
		userIDs[index] = meetingUser.UserID
	}
	users := []user{}
	err = findByUniqueKey(userIDs, &users)
	if err != nil {
		return m
	}
	ma["users"] = users
	return ma
}

type meetingMessage struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Text      string `gorm:"type:text;not null"`
	MeetingID uint   `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	ParentID  *uint
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
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	MeetingID uint `gorm:"not null"`
	UserID    uint `gorm:"not null"`
}

func (_ *meetingUser) TableName() string {
	return "meeting_user"
}

type project struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(128);not null"`
	TeamID    uint   `gorm:"not null"`
}

func (p project) public() interface{} {
	i, err := convert(p)
	if err != nil {
		return p
	}
	m, ok := i.(map[string]interface{})
	if !ok {
		return p
	}
	projectUsers := []projectUser{}
	err = find(map[string]interface{}{"projectID": p.ID}, &projectUsers)
	if err != nil {
		return p
	}
	userIDs := make([]uint, len(projectUsers))
	for index, projectUser := range projectUsers {
		userIDs[index] = projectUser.UserID
	}
	users := []user{}
	err = findByUniqueKey(userIDs, &users)
	if err != nil {
		return p
	}
	m["users"] = users
	return m
}

type projectUser struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ProjectID uint `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	RoleID    uint `gorm:"not null"`
}

func (_ *projectUser) TableName() string {
	return "project_user"
}

type projectUserRole struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string `gorm:"not null;unique"`
}

type request struct {
	Text string `gorm:"type:text;not null"`
}

type scan struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CommitSHA1 string `gorm:"type:char(40);not null"`
	UserID     uint   `gorm:"not null"`
	ProjectID  uint   `gorm:"not null"`
}

func (s scan) public() interface{} {
	i, err := convert(s)
	if err != nil {
		return s
	}
	m := i.(map[string]interface{})
	vulnerabilities := []vulnerability{}
	err = find(map[string]interface{}{"scanID": s.ID}, &vulnerabilities)
	if err != nil {
		m["vulnerabilities"] = []vulnerability{}
		return m
	}
	vp := make([]interface{}, len(vulnerabilities))
	for index, v := range vulnerabilities {
		p, err := public(v)
		if err != nil {
			continue
		}
		vp[index] = p
	}
	m["vulnerabilities"] = vp
	return m
}

type team struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(256);not null;unique;"`
}

func (t team) public() interface{} {
	i, err := convert(t)
	if err != nil {
		return t
	}
	m, ok := i.(map[string]interface{})
	if !ok {
		return t
	}
	projects := []project{}
	err = find(map[string]interface{}{"teamID": t.ID}, &projects)
	if err != nil {
		return t
	}
	m["projects"] = projects
	teamUsers := []teamUser{}
	err = find(map[string]interface{}{"teamID": t.ID}, &teamUsers)
	if err != nil {
		return t
	}
	users := make([]user, len(teamUsers))
	for index, teamUser := range teamUsers {
		u := user{}
		err = first(map[string]interface{}{"id": teamUser.UserID}, &u)
		if err != nil {
			return t
		}
		users[index] = u
	}
	m["users"] = users
	return m
}

type teamUser struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TeamID    uint `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	RoleID    uint `gorm:"not null"`
}

func (_ *teamUser) TableName() string {
	return "team_user"
}

func (t teamUser) public() interface{} {
	i, err := convert(t)
	if err != nil {
		return t
	}
	m, ok := i.(map[string]interface{})
	if !ok {
		return t
	}
	teamUserRole := teamUserRole{}
	err = first(map[string]interface{}{"id": t.RoleID}, &teamUserRole)
	if err != nil {
		return t
	}
	m["role"] = teamUserRole
	return m
}

type teamUserInvitationRequest struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Message       string `gorm:"type:varchar(512);"`
	TeamID        uint   `gorm:"not null"`
	InviterUserID uint   `gorm:"not null"`
	InviteeUserID uint   `gorm:"not null"`
	RoleID        uint   `gorm:"not null"`
}

func (t teamUserInvitationRequest) public() interface{} {
	i, err := convert(t)
	if err != nil {
		return t
	}
	m, ok := i.(map[string]interface{})
	if !ok {
		return t
	}
	inviterUser := user{}
	err = first(map[string]interface{}{"id": t.InviterUserID}, &inviterUser)
	if err != nil {
		return t
	}
	m["inviterUser"] = inviterUser
	inviteeUser := user{}
	err = first(map[string]interface{}{"id": t.InviteeUserID}, &inviteeUser)
	if err != nil {
		return t
	}
	m["inviteeUser"] = inviteeUser
	r := teamUserRole{}
	err = first(map[string]interface{}{"id": t.RoleID}, &r)
	if err != nil {
		return t
	}
	m["role"] = r
	teamUserInvitationRequestProjects := []teamUserInvitationRequestProject{}
	err = find(map[string]interface{}{"teamUserInvitationRequestID": t.ID}, &teamUserInvitationRequestProjects)
	if err != nil {
		return t
	}
	projects := make([]project, len(teamUserInvitationRequestProjects))
	for index, teamUserInvitationRequestProject := range teamUserInvitationRequestProjects {
		p := project{}
		err = first(map[string]interface{}{"id": teamUserInvitationRequestProject.ProjectID}, &p)
		if err != nil {
			return t
		}
		projects[index] = p
	}
	m["projects"] = projects
	team := team{}
	err = first(map[string]interface{}{"id": t.TeamID}, &team)
	if err != nil {
		return t
	}
	m["team"] = team
	return m
}

type teamUserInvitationRequestProject struct {
	TeamUserInvitationRequestID uint `gorm:"not null"`
	ProjectID                   uint `gorm:"not null"`
}

type teamUserRole struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string `gorm:"not null;unique"`
}

type test struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Steps      int    `gorm:"not null"`
	Branchname string `gorm:"type:varchar(256); not null;"`
	CommitSHA1 string `gorm:"type:char(40);not null;unique"`
	ProjectID  uint   `gorm:"not null"`
}

func (t test) public() interface{} {
	i, err := convert(t)
	if err != nil {
		return t
	}
	m := i.(map[string]interface{})
	testResults := []testResult{}
	err = find(map[string]interface{}{"testID": t.ID}, &testResults)
	if err != nil {
		return t
	}
	statusText := testStatusSuccessText
	if t.Steps != len(testResults) {
		statusText = testStatusRunningText
	}
	for _, testResult := range testResults {
		testStatus := testStatus{}
		err = first(map[string]interface{}{"id": testResult.TestStatusID}, &testStatus)
		if err != nil {
			break
		}
		if testStatus.Text == testStatusFailedText {
			statusText = testStatusFailedText
		}
	}
	rp := make([]interface{}, len(testResults))
	for index, testResult := range testResults {
		if statusText == testStatusSuccessText {
			status := testStatus{}
			err = first(map[string]interface{}{"id": testResult.TestStatusID}, &status)
			if err != nil {
				continue
			}
			statusText = status.Text
		}
		rp[index] = testResult.public()
	}
	m["color"] = testStatusColors[statusText]
	m["results"] = rp
	m["status"] = statusText
	return m
}

type testStatus struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Text      string `gorm:"type:varchar(7);unique"`
}

type testResult struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Command      string `gorm:"type:text;not null"`
	Output       string `gorm:"type:text;"`
	TestID       uint   `gorm:"not null"`
	TestStatusID uint   `gorm:"not null"`
	CompletedAt  *time.Time
}

func (t testResult) public() interface{} {
	i, err := convert(t)
	if err != nil {
		return t
	}
	m := i.(map[string]interface{})
	ts := testStatus{}
	err = first(map[string]interface{}{"id": t.TestStatusID}, &ts)
	if err != nil {
		return t
	}
	m["color"] = testStatusColors[ts.Text]
	m["status"] = ts.Text
	return m
}

type user struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Name             string `gorm:"type:varchar(32);not null;unique"`
	Password         string `gorm:"type:char(60);not null;"`
	ProfileImagePath string `gorm:"type:varchar(256);not null;"`
}

type vulnerability struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `gorm:"type:varchar(32);not null"`
	Description string `gorm:"type:varchar(128);not null"`
	Path        string `gorm:"type:varchar(256);not null"`
	Method      string `gorm:"type:varchar(8);not null"`
	Request     string `gorm:"type:text;not null"`
	Response    string `gorm:"type:text;not null"`
	ProjectID   uint   `gorm:"not null"`
	ScanID      uint   `gorm:"not null"`
}
