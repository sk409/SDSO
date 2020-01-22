package main

import (
	"bytes"
	"io/ioutil"
	"log"
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
		respondJSON(w, http.StatusOK, response(false))
		return
	}
	session, err := sessionManager.Provider.Get(sessionCookie.Value)
	if session == nil || err != nil {
		respondJSON(w, http.StatusOK, response(false))
		return
	}
	respondJSON(w, http.StatusOK, response(true))
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
	username := r.URL.Query().Get("username")
	projectname := r.URL.Query().Get("projectname")
	if emptyAny(username, projectname) {
		respond(w, http.StatusBadRequest)
		return
	}
	branches, err := gitRepositories.Branches(filepath.Join(username, projectname))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	branchnames := []string{}
	for _, branch := range branches {
		branchnames = append(branchnames, string(branch))
	}
	respondJSON(w, http.StatusOK, branchnames)
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
	respondJSON(w, http.StatusOK, branchProtectionRules)
}

func (b *branchProtectionRulesHandler) store(w http.ResponseWriter, r *http.Request) {
	branchProtectionRule := branchProtectionRule{}
	statusCode, err := store(r, &branchProtectionRule)
	if err != nil {
		respondError(w, statusCode, err)
		return
	}
	respondJSON(w, http.StatusOK, branchProtectionRule)
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
	username := r.URL.Query().Get("username")
	projectname := r.URL.Query().Get("projectname")
	branchname := r.URL.Query().Get("branchname")
	if emptyAny(username, projectname, branchname) {
		respond(w, http.StatusBadRequest)
		return
	}
	commitsByte, err := gitRepositories.Log(filepath.Join(username, projectname), branchname, "--pretty=format:%h %cd %s")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
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
	respondJSON(w, http.StatusOK, commits)
}

func (c *commitsHandler) show(w http.ResponseWriter, r *http.Request, sha1 string) {
	userName := r.URL.Query().Get("userName")
	projectName := r.URL.Query().Get("projectName")
	if emptyAny(userName, projectName) {
		respond(w, http.StatusBadRequest)
		return
	}
	commits, err := gitRepositories.Log(filepath.Join(userName, projectName), "--pretty=oneline")
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
	diff, err := gitRepositories.Show(filepath.Join(userName, projectName), sha1)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commit.Diff = string(diff)
	respondJSON(w, http.StatusOK, commit)
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
	username := r.URL.Query().Get("username")
	projectname := r.URL.Query().Get("projectname")
	treeIsh := r.URL.Query().Get("treeIsh")
	path := r.URL.Query().Get("path")
	if emptyAny(username, projectname, treeIsh) {
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
	output, err := gitRepositories.LsTree(filepath.Join(username, projectname), treeIsh, args...)
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
	respondJSON(w, http.StatusOK, files)
}

func (f *filesHandler) text(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	projectname := r.URL.Query().Get("projectname")
	revision := r.URL.Query().Get("revision")
	path := r.URL.Query().Get("path")
	if emptyAny(username, projectname, revision, path) {
		respond(w, http.StatusBadRequest)
		return
	}
	output, err := gitRepositories.LsTree(filepath.Join(username, projectname), revision, path)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	regex := regexp.MustCompile("[0-9]+ [a-z]+ ([0-9a-z]+)")
	matches := regex.FindSubmatch(output)
	log.Println(string(output))
	log.Println(matches)
	if len(matches) != 2 {
		respond(w, http.StatusBadRequest)
		return
	}
	sha1 := string(matches[1])
	text, err := gitRepositories.CatFile(filepath.Join(username, projectname), sha1, "-p")
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

func (g *gitHandler) receivePack(w http.ResponseWriter, r *http.Request, userName, projectName string) {
	user := user{}
	db.Where("name = ?", userName).First(&user)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	project := project{}
	db.Where("name = ? AND user_id = ?", projectName, user.ID).First(&project)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if project.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	branchProtectionRules := []branchProtectionRule{}
	db.Where("project_id = ?", project.ID).Find(&branchProtectionRules)
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
		repositoryPath := filepath.Join(gitRepositories.RootDirectoryPath, userName, projectName)
		tmpRepositoryPath := filepath.Join(gitTmpRepositories.RootDirectoryPath, userName, projectName)
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
		defer os.RemoveAll(filepath.Join(gitTmpRepositories.RootDirectoryPath, userName))
		err = gofile.Copy(repositoryPath, tmpRepositoryPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = gitTmpRepositories.RPC(filepath.Join(userName, projectName), gogit.RPCReceivePack, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		succeeded, err := runTest(userName, projectName, tmpRepositoryPath, branchName, commitSHA1)
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
		clonePath := filepath.Join(gitRepositories.RootDirectoryPath, userName, projectName)
		go runTest(userName, projectName, clonePath, branchName, commitSHA1)
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
	respondJSON(w, http.StatusOK, response)
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

type projectsHandler struct {
}

func (p *projectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (p *projectsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	projects := []project{}
	err := fetch(r, &projects)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, projects)
}

func (p *projectsHandler) store(w http.ResponseWriter, r *http.Request) {
	project := project{}
	statusCode, err := store(r, &project)
	if err != nil {
		respondError(w, statusCode, err)
		return
	}
	respondJSON(w, http.StatusOK, project)
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
		Name:     username,
		Password: string(hashedPassword),
	}
	count := 0
	db.Model(&user{}).Where("name = ?", u.Name).Count(&count)
	if count != 0 {
		respondJSON(w, http.StatusOK, response(nil, false))
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
	respondJSON(w, http.StatusOK, response(&u, true))
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
	scans := []scan{}
	err := fetch(r, &scans)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, scans)
}

func (s *scansHandler) store(w http.ResponseWriter, r *http.Request) {
	userName := r.PostFormValue("userName")
	projectName := r.PostFormValue("projectName")
	if emptyAny(userName, projectName) {
		respond(w, http.StatusBadRequest)
		return
	}
	u := user{}
	statusCode, err := first(map[string]interface{}{"name": userName}, &u)
	if err != nil {
		respondError(w, statusCode, err)
		return
	}
	p := project{}
	statusCode, err = first(map[string]interface{}{"name": projectName, "user_id": u.ID}, &p)
	if err != nil {
		respondError(w, statusCode, err)
		return
	}
	scan := scan{
		UserID:    u.ID,
		ProjectID: p.ID,
	}
	db.Save(&scan)
	if db.Error != nil {
		respondError(w, http.StatusInternalServerError, db.Error)
		return
	}
	respondJSON(w, http.StatusOK, scan)
}

type testsHandler struct {
}

func (t *testsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/tests"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			t.fetch(w, r)
			return
		case path.Join(base, "branch"):
			t.branch(w, r)
			return
		case path.Join(base, "socket"):
			t.socket(w, r)
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
	respondJSON(w, http.StatusOK, tests)
}

func (t *testsHandler) branch(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("userName")
	projectName := r.URL.Query().Get("projectName")
	projectID := r.URL.Query().Get("projectID")
	branchName := r.URL.Query().Get("branchName")
	if emptyAny(userName, projectName, projectID, branchName) {
		respond(w, http.StatusBadRequest)
		return
	}
	tests := []test{}
	db.Where("project_id = ?", projectID).Find(&tests)
	if db.Error != nil {
		respondError(w, http.StatusInternalServerError, db.Error)
		return
	}
	if len(tests) == 0 {
		respondJSON(w, http.StatusOK, []test{})
		return
	}
	commitSHA1sByte, err := gitRepositories.RevList(filepath.Join(userName, projectName), branchName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	commitSHA1s := strings.Split(string(commitSHA1sByte), "\n")
	response := []test{}
	for _, test := range tests {
		for _, commitSHA1 := range commitSHA1s {
			if test.CommitSHA1 == commitSHA1 {
				response = append(response, test)
			}
		}
	}
	respondJSON(w, http.StatusOK, response)
}

func (t *testsHandler) socket(w http.ResponseWriter, r *http.Request) {
	u, err := authenticatedUser(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	socket, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	websocketsTest[u.ID] = socket
}

type testResultsHandler struct {
}

func (t *testResultsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/test_results/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			t.fetch(w, r)
			return
		case path.Join(base, "socket"):
			t.socket(w, r)
			return
		}
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
	respondJSON(w, http.StatusOK, testResults)
}

func (t *testResultsHandler) socket(w http.ResponseWriter, r *http.Request) {
	u, err := authenticatedUser(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	socket, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	websocketsTestResult[u.ID] = socket
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
	respondJSON(w, http.StatusOK, testStatuses)
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
	respondJSON(w, http.StatusOK, user)
}

type usersHandler struct {
}

func (u *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u.fetch(w, r)
		return
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
	respondJSON(w, http.StatusOK, users)
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
	respondJSON(w, http.StatusOK, vulnerabilities)
}

func (v *vulnerabilitiesHandler) store(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	path := r.PostFormValue("path")
	method := r.PostFormValue("method")
	request := r.PostFormValue("request")
	response := r.PostFormValue("response")
	scanID := r.PostFormValue("scanID")
	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")
	projectName := r.PostFormValue("projectName")
	if emptyAny(name, description, path, method, request, response, scanID, userName, password, projectName) {
		respond(w, http.StatusBadRequest)
		return
	}
	u := user{}
	statusCode, err := first(map[string]interface{}{"name": userName}, &u)
	if err != nil {
		respondError(w, statusCode, err)
		return
	}
	p := project{}
	statusCode, err = first(map[string]interface{}{"name": projectName, "user_id": u.ID}, &p)
	if err != nil {
		respondError(w, statusCode, err)
		return
	}
	s := scan{}
	statusCode, err = first(map[string]interface{}{"id": scanID}, &s)
	if err != nil {
		respondError(w, statusCode, err)
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
	respondJSON(w, http.StatusOK, vulnerability)
}
