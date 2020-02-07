package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/sk409/gofile"
	"github.com/sk409/gogit"
)

type authHandler struct {
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.auth(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (a *authHandler) auth(w http.ResponseWriter, r *http.Request) {
	response := func(authenticated bool) map[string]bool {
		return map[string]bool{"authenticated": authenticated}
	}
	sessionCookie, err := r.Cookie("sessionID")
	if err != nil {
		_, err = respondJSON(w, http.StatusOK, response(false))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	session, err := sessionManager.Provider.Get(sessionCookie.Value)
	if session == nil || err != nil {
		_, err = respondJSON(w, http.StatusOK, response(false))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	_, err = respondJSON(w, http.StatusOK, response(true))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type branchesHandler struct {
}

func (b *branchesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b.fetch(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (b *branchesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	teamname := r.URL.Query().Get("teamname")
	projectname := r.URL.Query().Get("projectname")
	if emptyAny(teamname, projectname) {
		respond(w, http.StatusBadRequest)
		return
	}
	branches, err := gitRepositories.Branches(filepath.Join(teamname, projectname))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	branchnames := []string{}
	for _, branch := range branches {
		branchnames = append(branchnames, string(branch))
	}
	_, err = respondJSON(w, http.StatusOK, branchnames)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type branchProtectionRulesHandler struct {
}

func (b *branchProtectionRulesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b.fetch(w, r)
		return
	case http.MethodPost:
		b.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (b *branchProtectionRulesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	branchProtectionRules := []branchProtectionRule{}
	err := fetch(r, &branchProtectionRules)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, branchProtectionRules)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (b *branchProtectionRulesHandler) store(w http.ResponseWriter, r *http.Request) {
	branchProtectionRule := branchProtectionRule{}
	err := store(r, &branchProtectionRule)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, branchProtectionRule)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type commitsHandler struct {
}

func (c *commitsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/commits"
	switch r.Method {
	case http.MethodGet:
		sha1, ok := routeWithID(r)
		if ok {
			c.show(w, r, sha1)
			return
		} else if r.URL.Path == base+"/" {
			c.fetch(w, r)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (c *commitsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	teamname := r.URL.Query().Get("teamname")
	projectname := r.URL.Query().Get("projectname")
	branchname := r.URL.Query().Get("branchname")
	if emptyAny(teamname, projectname, branchname) {
		respond(w, http.StatusBadRequest)
		return
	}
	commitsByte, err := gitRepositories.Log(filepath.Join(teamname, projectname), branchname, "--pretty=format:%h %cd %s")
	if err != nil {
		_, err = respondJSON(w, http.StatusOK, []commit{})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	lines := strings.Split(string(commitsByte), "\n")
	commits := []commit{}
	for _, line := range lines {
		components := strings.Split(line, " ")
		if len(components) < 8 {
			continue
		}
		sha1 := components[0]
		date := strings.Join(components[1:7], " ")
		message := strings.Join(components[7:], " ")
		commits = append(commits, commit{
			Branchname: branchname,
			SHA1:       sha1,
			Date:       date,
			Message:    message,
		})
	}
	_, err = respondJSON(w, http.StatusOK, commits)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (c *commitsHandler) show(w http.ResponseWriter, r *http.Request, sha1 string) {
	teamname := r.URL.Query().Get("teamname")
	projectname := r.URL.Query().Get("projectname")
	if emptyAny(teamname, projectname) {
		respond(w, http.StatusBadRequest)
		return
	}
	commits, err := gitRepositories.Log(filepath.Join(teamname, projectname), "--pretty=oneline")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	lines := strings.Split(string(commits), "\n")
	commit := commit{SHA1: sha1}
	for _, line := range lines {
		if strings.Contains(line, sha1) {
			message := strings.Split(line, " ")[1]
			commit.Message = message
			break
		}
	}
	diff, err := gitRepositories.Show(filepath.Join(teamname, projectname), sha1)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	commit.Diff = string(diff)
	_, err = respondJSON(w, http.StatusOK, commit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type dastVulnerabilityMessagesHandler struct {
}

func (d *dastVulnerabilityMessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/dast_vulnerability_messages/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			d.fetch(w, r)
			return
		case base + "socket":
			d.socket(w, r)
			return
		}
	case http.MethodPost:
		d.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (d *dastVulnerabilityMessagesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	dastVulnerabilityMessages := []dastVulnerabilityMessage{}
	err := fetch(r, &dastVulnerabilityMessages)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, dastVulnerabilityMessages)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (d *dastVulnerabilityMessagesHandler) socket(w http.ResponseWriter, r *http.Request) {
	websock(w, r, &websocketsDastVulnerabilityMessage)
}

func (d *dastVulnerabilityMessagesHandler) store(w http.ResponseWriter, r *http.Request) {
	dastVulnerabilityMessage := dastVulnerabilityMessage{}
	err := store(r, &dastVulnerabilityMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, dastVulnerabilityMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	v := vulnerability{}
	err = first(map[string]interface{}{"id": dastVulnerabilityMessage.VulnerabilityID}, &v)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	projectUsers := []projectUser{}
	err = find(map[string]interface{}{"project_id": v.ProjectID}, &projectUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := public(dastVulnerabilityMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	for _, projectUser := range projectUsers {
		socket, ok := websocketsDastVulnerabilityMessage[projectUser.UserID]
		if !ok {
			continue
		}
		socket.WriteJSON(p)
	}
}

type filesHandler struct {
}

func (f *filesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/files"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base + "/":
			f.fetch(w, r)
			return
		case path.Join(base, "text"):
			f.text(w, r)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (f *filesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	teamname := r.URL.Query().Get("teamname")
	projectname := r.URL.Query().Get("projectname")
	revision := r.URL.Query().Get("revision")
	path := r.URL.Query().Get("path")
	if emptyAny(teamname, projectname, revision) {
		respond(w, http.StatusBadRequest)
		return
	}
	args := []string{}
	if !emptyAny(path) {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		args = append(args, path)
	}
	output, err := gitRepositories.LsTree(filepath.Join(teamname, projectname), revision, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	regex := regexp.MustCompile("[0-9]+ ([a-z]+) [a-z0-9]+\\t(.+)")
	files := []map[string]interface{}{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		matches := regex.FindAllStringSubmatch(line, -1)
		if len(matches) == 0 {
			continue
		}
		match := matches[0]
		if len(match) != 3 {
			continue
		}
		kind := match[1]
		filename := strings.TrimPrefix(match[2], path)
		file := map[string]interface{}{}
		file["path"] = filepath.Join(path, filename)
		file["name"] = filename
		file["isDirectory"] = kind == "tree"
		files = append(files, file)
	}
	_, err = respondJSON(w, http.StatusOK, files)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (f *filesHandler) text(w http.ResponseWriter, r *http.Request) {
	teamname := r.URL.Query().Get("teamname")
	projectname := r.URL.Query().Get("projectname")
	revision := r.URL.Query().Get("revision")
	path := r.URL.Query().Get("path")
	if emptyAny(teamname, projectname, revision, path) {
		respond(w, http.StatusBadRequest)
		return
	}
	output, err := gitRepositories.LsTree(filepath.Join(teamname, projectname), revision, path)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	regex := regexp.MustCompile("[0-9]+ [a-z]+ ([0-9a-z]+)")
	matches := regex.FindSubmatch(output)
	if len(matches) != 2 {
		respond(w, http.StatusBadRequest)
		return
	}
	sha1 := string(matches[1])
	text, err := gitRepositories.CatFile(filepath.Join(teamname, projectname), sha1, "-p")
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondMessage(w, http.StatusOK, string(text))
}

type gitHandler struct {
}

func (g *gitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	components := strings.Split(r.URL.Path, "/")
	if len(components) < 3 {
		respond(w, http.StatusNotFound)
		return
	}
	userName := components[1]
	projectName := components[2]
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/"+path.Join(userName, projectName, "info", "refs") {
			g.refs(w, r)
			return
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/" + path.Join(userName, projectName, "git-receive-pack"):
			g.receivePack(w, r, userName, projectName)
			return
		case "/" + path.Join(userName, projectName, "git-upload-pack"):
			g.uploadPack(w, r)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (g *gitHandler) refs(w http.ResponseWriter, r *http.Request) {
	gitServer.ServeHTTP(w, r)
}

func (g *gitHandler) receivePack(w http.ResponseWriter, r *http.Request, teamname, projectname string) {
	t := team{}
	err := first(map[string]interface{}{"name": teamname}, &t)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if t.ID == 0 {
		respondError(w, http.StatusInternalServerError, errBadRequest)
		return
	}
	p := project{}
	err = first(map[string]interface{}{"name": projectname, "team_id": t.ID}, &p)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if p.ID == 0 {
		respondError(w, http.StatusInternalServerError, errBadRequest)
		return
	}
	branchProtectionRules := []branchProtectionRule{}
	db.Where("project_id = ?", p.ID).Find(&branchProtectionRules)
	branchName, commitSHA1, err := getBranchNameAndCommitSHA1(r)
	if err != nil {
		return
	}
	protection := false
	for _, branchProtectionRule := range branchProtectionRules {
		if branchName == branchProtectionRule.BranchName {
			protection = true
			break
		}
	}
	if protection {
		body, err := gogit.GetReadCloser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		repositoryPath := filepath.Join(gitRepositories.RootDirectoryPath, teamname, projectname)
		tmpRepositoryPath := filepath.Join(gitTmpRepositories.RootDirectoryPath, teamname, projectname)
		err = os.RemoveAll(tmpRepositoryPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = os.MkdirAll(tmpRepositoryPath, 0755)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer os.RemoveAll(filepath.Join(gitTmpRepositories.RootDirectoryPath, teamname))
		err = gofile.Copy(repositoryPath, tmpRepositoryPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = gitTmpRepositories.RPC(filepath.Join(teamname, projectname), gogit.RPCReceivePack, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tester := tester{}
		succeeded, err := tester.run(teamname, projectname, tmpRepositoryPath, branchName, commitSHA1)
		if !succeeded || err != nil {
			// log.Println("****************")
			// log.Println(err)
			//
			// TODO: ちゃんとしたレスポンスを返す
			w.WriteHeader(http.StatusInternalServerError)
			return
			//
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		gitServer.ServeHTTP(w, r)
	} else {
		gitServer.ServeHTTP(w, r)
		clonePath := filepath.Join(gitRepositories.RootDirectoryPath, teamname, projectname)
		tester := tester{}
		go tester.run(teamname, projectname, clonePath, branchName, commitSHA1)
	}
}

func (g *gitHandler) uploadPack(w http.ResponseWriter, r *http.Request) {
	gitServer.ServeHTTP(w, r)
}

type loginHandler struct {
}

func (l *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		l.login(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (l *loginHandler) login(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if emptyAny(username, password) {
		respond(w, http.StatusBadRequest)
		return
	}
	u, err := login(w, username, password)
	response := map[string]interface{}{
		"ok":   err == nil,
		"user": u,
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type logoutHandler struct {
}

func (l *logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		l.logout(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (l *logoutHandler) logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	sessionID := sessionCookie.Value
	_, err = sessionManager.Provider.Get(sessionID)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	sessionManager.Provider.Stop(sessionID)
	http.SetCookie(w, &http.Cookie{
		Name:   cookieNameSessionID,
		MaxAge: -1,
	})
}

type meetingsHandler struct {
}

func (m *meetingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/meetings/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			m.fetch(w, r)
			return
		case base + "ids":
			m.ids(w, r)
			return
		}
	case http.MethodPost:
		m.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (m *meetingsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	meetings := []meeting{}
	err := fetch(r, &meetings)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, meetings)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (m *meetingsHandler) ids(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["ids[]"]
	if emptyAny(ids) {
		_, err := respondJSON(w, http.StatusOK, []meeting{})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	meetings := []meeting{}
	err := findByUniqueKey(ids, &meetings)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, meetings)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (m *meetingsHandler) store(w http.ResponseWriter, r *http.Request) {
	meeting := meeting{}
	err := store(r, &meeting)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, meeting)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type meetingMessagesHandler struct {
}

func (m *meetingMessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/meeting_messages/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			m.fetch(w, r)
			return
		case base + "socket":
			m.socket(w, r)
			return
		}
	case http.MethodPost:
		m.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (m *meetingMessagesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	meetingMessages := []meetingMessage{}
	err := fetch(r, &meetingMessages)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, meetingMessages)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (m *meetingMessagesHandler) socket(w http.ResponseWriter, r *http.Request) {
	websock(w, r, &websocketsMeetingMessage)
}

func (m *meetingMessagesHandler) store(w http.ResponseWriter, r *http.Request) {
	meetingMessage := meetingMessage{}
	err := store(r, &meetingMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, meetingMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	meetingUsers := []meetingUser{}
	err = find(map[string]interface{}{"meetingID": meetingMessage.MeetingID}, &meetingUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	for _, meetingUser := range meetingUsers {
		socket, ok := websocketsMeetingMessage[meetingUser.UserID]
		if !ok {
			continue
		}
		p, err := public(meetingMessage)
		if err != nil {
			continue
		}
		socket.WriteJSON(p)
	}
}

type meetingUsersHandler struct {
}

func (m *meetingUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.fetch(w, r)
		return
	case http.MethodPost:
		m.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (m *meetingUsersHandler) fetch(w http.ResponseWriter, r *http.Request) {
	meetingUsers := []meetingUser{}
	err := fetch(r, &meetingUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, meetingUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (m *meetingUsersHandler) store(w http.ResponseWriter, r *http.Request) {
	meetingUser := meetingUser{}
	err := store(r, &meetingUser)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, meetingUser)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type projectsHandler struct {
}

func (p *projectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/projects/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			p.fetch(w, r)
			return
		case base + "ids":
			p.ids(w, r)
			return
		}
	case http.MethodPost:
		p.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (p *projectsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	projects := []project{}
	err := fetch(r, &projects)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, projects)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (p *projectsHandler) ids(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := r.Form["ids[]"]
	if emptyAny(ids) {
		_, err := respondJSON(w, http.StatusOK, []project{})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	projects := []project{}
	err := findByUniqueKey(ids, &projects)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, projects)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (p *projectsHandler) store(w http.ResponseWriter, r *http.Request) {
	project := project{}
	err := store(r, &project)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	t := team{}
	err = first(map[string]interface{}{"id": project.TeamID}, &t)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	err = os.MkdirAll(filepath.Join(pathRepositories, t.Name, project.Name), 0755)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	err = gitRepositories.InitBare(filepath.Join(t.Name, project.Name))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, project)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type registerHandler struct {
}

func (h *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.register(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (h *registerHandler) register(w http.ResponseWriter, r *http.Request) {
	response := func(u *user, ok bool) map[string]interface{} {
		return map[string]interface{}{
			"ok":   ok,
			"user": u,
		}
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if emptyAny(username, password) {
		respond(w, http.StatusNotFound)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	u := user{
		Name:             username,
		Password:         string(hashedPassword),
		ProfileImagePath: pathNoImage,
	}
	count := 0
	db.Model(&user{}).Where("name = ?", u.Name).Count(&count)
	if count != 0 {
		_, err = respondJSON(w, http.StatusOK, response(nil, false))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	db.Save(&u)
	if db.Error != nil {
		respondError(w, http.StatusInternalServerError, db.Error)
		return
	}
	_, err = login(w, username, password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, response(&u, true))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type projectUsersHandler struct {
}

func (p *projectUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.fetch(w, r)
		return
	case http.MethodPost:
		p.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (p *projectUsersHandler) fetch(w http.ResponseWriter, r *http.Request) {
	projectUsers := []projectUser{}
	err := fetch(r, &projectUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, projectUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (p *projectUsersHandler) store(w http.ResponseWriter, r *http.Request) {
	projectID := r.PostFormValue("projectId")
	role := r.PostFormValue("role")
	userID := r.PostFormValue("userId")
	if emptyAny(projectID, role, userID) {
		respond(w, http.StatusBadRequest)
		return
	}
	projectUserRole := projectUserRole{}
	err := first(map[string]interface{}{"role": role}, &projectUserRole)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	projectUser := projectUser{}
	err = save(map[string]interface{}{"projectId": projectID, "roleId": projectUserRole.ID, "userId": userID}, &projectUser)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusInternalServerError, projectUser)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type repositoriesHandler struct {
}

func (h *repositoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (h *repositoriesHandler) store(w http.ResponseWriter, r *http.Request) {
	userName := r.PostFormValue("userName")
	projectName := r.PostFormValue("projectName")
	if emptyAny(userName, projectName) {
		respond(w, http.StatusBadRequest)
		return
	}
	repositoryPath := filepath.Join(userName, projectName)
	repositoriDirectoryPath := filepath.Join(gitRepositories.RootDirectoryPath, repositoryPath)
	exist := gofile.IsExist(repositoriDirectoryPath)
	if exist {
		respond(w, http.StatusBadRequest)
		return
	}
	err := os.MkdirAll(repositoriDirectoryPath, 0755)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	err = gitRepositories.InitBare(repositoryPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type scansHandler struct {
}

func (s *scansHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.fetch(w, r)
		return
	case http.MethodPost:
		s.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (s *scansHandler) fetch(w http.ResponseWriter, r *http.Request) {
	projectname := r.URL.Query().Get("projectname")
	revision := r.URL.Query().Get("revision")
	teamname := r.URL.Query().Get("teamname")
	if emptyAny(projectname, revision, teamname) {
		respond(w, http.StatusBadRequest)
		return
	}
	t := team{}
	err := first(map[string]interface{}{"name": teamname}, &t)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p := project{}
	err = first(map[string]interface{}{"name": projectname, "team_id": t.ID}, &p)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	scans := []scan{}
	err = find(map[string]interface{}{"projectId": p.ID}, &scans)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	revisionsByte, err := gitRepositories.RevList(filepath.Join(teamname, projectname), revision)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	revisions := strings.Split(string(revisionsByte), "\n")
	response := []scan{}
	for _, scan := range scans {
		for _, revision := range revisions {
			if scan.CommitSHA1 == revision {
				response = append(response, scan)
				break
			}
		}
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *scansHandler) store(w http.ResponseWriter, r *http.Request) {
	commitSHA1 := r.PostFormValue("commitSHA1")
	username := r.PostFormValue("username")
	teamname := r.PostFormValue("teamname")
	projectname := r.PostFormValue("projectname")
	if emptyAny(commitSHA1, username, teamname, projectname) {
		respond(w, http.StatusBadRequest)
		return
	}
	u := user{}
	err := first(map[string]interface{}{"name": username}, &u)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	t := team{}
	err = first(map[string]interface{}{"name": teamname}, &t)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p := project{}
	err = first(map[string]interface{}{"name": projectname, "team_id": t.ID}, &p)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	scan := scan{
		CommitSHA1: commitSHA1,
		UserID:     u.ID,
		ProjectID:  p.ID,
	}
	db.Save(&scan)
	if db.Error != nil {
		respondError(w, http.StatusInternalServerError, db.Error)
		return
	}
	_, err = respondJSON(w, http.StatusOK, scan)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type teamsHandler struct {
}

func (t *teamsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/teams/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			t.fetch(w, r)
			return
		case base + "ids":
			t.ids(w, r)
			return
		}
	case http.MethodPost:
		t.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (t *teamsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	teams := []team{}
	err := fetch(r, &teams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamsHandler) ids(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["ids[]"]
	if emptyAny(ids) {
		_, err := respondJSON(w, http.StatusOK, []team{})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	teams := []team{}
	err := findByUniqueKey(ids, &teams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamsHandler) store(w http.ResponseWriter, r *http.Request) {
	team := team{}
	err := store(r, &team)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, team)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type teamUsersHandler struct {
}

func (t *teamUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.fetch(w, r)
		return
	case http.MethodPost:
		t.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (t *teamUsersHandler) fetch(w http.ResponseWriter, r *http.Request) {
	teamUsers := []teamUser{}
	err := fetch(r, &teamUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUsersHandler) store(w http.ResponseWriter, r *http.Request) {
	teamID := r.PostFormValue("teamId")
	role := r.PostFormValue("role")
	userID := r.PostFormValue("userId")
	if emptyAny(teamID, role, userID) {
		respond(w, http.StatusBadRequest)
		return
	}
	teamUserRole := teamUserRole{}
	err := first(map[string]interface{}{"role": role}, &teamUserRole)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	teamUser := teamUser{}
	err = save(map[string]interface{}{"teamID": teamID, "roleID": teamUserRole.ID, "userID": userID}, &teamUser)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusInternalServerError, teamUser)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type teamUserInvitationRequestProjectsHandler struct {
}

func (t *teamUserInvitationRequestProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.fetch(w, r)
		return
	case http.MethodPost:
		t.store(w, r)
		return
	case http.MethodDelete:
		t.destroy(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (t *teamUserInvitationRequestProjectsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	teamUserInvitationRequests := []teamUserInvitationRequest{}
	err := fetch(r, &teamUserInvitationRequests)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUserInvitationRequests)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUserInvitationRequestProjectsHandler) store(w http.ResponseWriter, r *http.Request) {
	teamUserInvitationRequestProject := teamUserInvitationRequestProject{}
	err := store(r, &teamUserInvitationRequestProject)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUserInvitationRequestProject)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUserInvitationRequestProjectsHandler) destroy(w http.ResponseWriter, r *http.Request) {
	teamUserInvitationRequestProject := []teamUserInvitationRequestProject{}
	err := destroy(r, &teamUserInvitationRequestProject)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUserInvitationRequestProject)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type teamUserInvitationRequestsHandler struct {
}

func (t *teamUserInvitationRequestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.fetch(w, r)
		return
	case http.MethodPost:
		t.store(w, r)
		return
	case http.MethodDelete:
		t.destroy(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (t *teamUserInvitationRequestsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	teamUserInvitationRequests := []teamUserInvitationRequest{}
	err := fetch(r, &teamUserInvitationRequests)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUserInvitationRequests)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUserInvitationRequestsHandler) store(w http.ResponseWriter, r *http.Request) {
	inviterUserID := r.PostFormValue("inviterUserId")
	inviteeUserID := r.PostFormValue("inviteeUserId")
	message := r.PostFormValue("message")
	roleText := r.PostFormValue("role")
	teamID := r.PostFormValue("teamId")
	if emptyAny(inviterUserID, inviteeUserID, message, roleText, teamID) {
		respond(w, http.StatusBadRequest)
		return
	}
	role := teamUserRole{}
	err := first(map[string]interface{}{"role": roleText}, &role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	teamUserInvitationRequest := teamUserInvitationRequest{}
	err = save(map[string]interface{}{"message": message, "teamID": teamID, "inviterUserID": inviterUserID, "inviteeUserID": inviteeUserID, "roleID": role.ID}, &teamUserInvitationRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUserInvitationRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUserInvitationRequestsHandler) destroy(w http.ResponseWriter, r *http.Request) {
	teamUserInvitationRequest := teamUserInvitationRequest{}
	err := destroy(r, &teamUserInvitationRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUserInvitationRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type testsHandler struct {
}

func (t *testsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/tests/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			t.fetch(w, r)
			return
		case base + "socket":
			t.socket(w, r)
			return
		case base + "revision":
			t.revision(w, r)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (t *testsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	tests := []test{}
	err := fetch(r, &tests)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, tests)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *testsHandler) socket(w http.ResponseWriter, r *http.Request) {
	websock(w, r, &websocketsTest)
}

func (t *testsHandler) revision(w http.ResponseWriter, r *http.Request) {
	teamname := r.URL.Query().Get("teamname")
	projectname := r.URL.Query().Get("projectname")
	revision := r.URL.Query().Get("revision")
	if emptyAny(teamname, projectname, revision) {
		respond(w, http.StatusBadRequest)
		return
	}
	team := team{}
	err := first(map[string]interface{}{"name": teamname}, &team)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p := project{}
	err = first(map[string]interface{}{"name": projectname, "team_id": team.ID}, &p)
	tests := []test{}
	db.Where("project_id = ?", p.ID).Order("created_at DESC").Find(&tests)
	if db.Error != nil {
		respondError(w, http.StatusInternalServerError, db.Error)
		return
	}
	if len(tests) == 0 {
		_, err = respondJSON(w, http.StatusOK, []test{})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	revisionsByte, err := gitRepositories.RevList(filepath.Join(teamname, projectname), revision)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	revisions := strings.Split(string(revisionsByte), "\n")
	response := []test{}
	for _, test := range tests {
		for _, revision := range revisions {
			if test.CommitSHA1 == revision {
				response = append(response, test)
				break
			}
		}
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type testMessagesHandler struct {
}

func (t *testMessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/test_messages/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			t.fetch(w, r)
			return
		case base + "socket":
			t.socket(w, r)
			return
		}
	case http.MethodPost:
		t.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (t *testMessagesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	testMessages := []testMessage{}
	err := fetch(r, &testMessages)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, testMessages)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *testMessagesHandler) socket(w http.ResponseWriter, r *http.Request) {
	websock(w, r, &websocketsTestMessage)
}

func (t *testMessagesHandler) store(w http.ResponseWriter, r *http.Request) {
	testMessage := testMessage{}
	err := store(r, &testMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, testMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	test := test{}
	err = first(map[string]interface{}{"id": testMessage.TestID}, &test)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	projectUsers := []projectUser{}
	err = find(map[string]interface{}{"project_id": test.ProjectID}, &projectUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := public(testMessage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	for _, projectUser := range projectUsers {
		socket, ok := websocketsTestMessage[projectUser.UserID]
		if !ok {
			continue
		}
		socket.WriteJSON(p)
	}
}

type testResultsHandler struct {
}

func (t *testResultsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.fetch(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (t *testResultsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	testResults := []testResult{}
	err := fetch(r, &testResults)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, testResults)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type testStatusesHandler struct {
}

func (t *testStatusesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.fetch(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (t *testStatusesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	testStatuses := []testStatus{}
	err := fetch(r, &testStatuses)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, testStatuses)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type userHandler struct {
}

func (u *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u.fetch(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (u *userHandler) fetch(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	session, err := sessionManager.Provider.Get(sessionCookie.Value)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	userID, err := session.Uint(sessionStoreNameUserID)
	user := user{}
	user.ID = userID
	db.Find(&user)
	if db.Error != nil {
		respondError(w, http.StatusInternalServerError, db.Error)
	}
	_, err = respondJSON(w, http.StatusOK, user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type usersHandler struct {
}

func (u *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/users/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			u.fetch(w, r)
			return
		case base + "ids":
			u.ids(w, r)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (u *usersHandler) fetch(w http.ResponseWriter, r *http.Request) {
	users := []user{}
	err := fetch(r, &users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (u *usersHandler) ids(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["ids[]"]
	if emptyAny(ids) {
		_, err := respondJSON(w, http.StatusOK, []user{})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	users := []user{}
	err := findByUniqueKey(ids, &users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

type vulnerabilitiesHandler struct {
}

func (v *vulnerabilitiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		v.fetch(w, r)
		return
	case http.MethodPost:
		v.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (v *vulnerabilitiesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	vulnerabilities := []vulnerability{}
	err := fetch(r, &vulnerabilities)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, vulnerabilities)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (v *vulnerabilitiesHandler) store(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	path := r.PostFormValue("path")
	method := r.PostFormValue("method")
	request := r.PostFormValue("request")
	response := r.PostFormValue("response")
	scanID := r.PostFormValue("scanID")
	teamname := r.PostFormValue("teamname")
	projectname := r.PostFormValue("projectname")
	if emptyAny(name, description, path, method, request, response, scanID, teamname, projectname) {
		respond(w, http.StatusBadRequest)
		return
	}
	t := team{}
	err := first(map[string]interface{}{"name": teamname}, &t)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p := project{}
	err = first(map[string]interface{}{"name": projectname, "team_id": t.ID}, &p)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	s := scan{}
	err = first(map[string]interface{}{"id": scanID}, &s)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
	vulnerability := vulnerability{
		Name:        name,
		Description: description,
		Path:        path,
		Method:      method,
		Request:     request,
		Response:    response,
		ProjectID:   p.ID,
		ScanID:      s.ID,
	}
	db.Save(&vulnerability)
	if db.Error != nil {
		respondError(w, http.StatusInternalServerError, db.Error)
		return
	}
	_, err = respondJSON(w, http.StatusOK, vulnerability)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}
