package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
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
	respond(w, http.StatusMethodNotAllowed)
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
	respond(w, http.StatusMethodNotAllowed)
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
	branchProtectionRules, err := branchProtectionRuleRepository.find(makeQuery(r, branchProtectionRule{}, true), loadAllRelation)
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
	projectID := r.PostFormValue("projectId")
	if emptyAny(projectID) {
		respond(w, http.StatusBadRequest)
		return
	}
	p, err := projectRepository.first(map[string]interface{}{"id": projectID}, projectRelationUsers)
	ok, err := checkPermissionWithRequest(r, p.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	branchProtectionRule, err := branchProtectionRuleRepository.save(makeQuery(r, branchProtectionRule{}, false))
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
		respond(w, http.StatusNotFound)
		return
	}
	respond(w, http.StatusMethodNotAllowed)
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
		case base + "count":
			d.count(w, r)
			return
		case base + "range":
			d.fetchRange(w, r)
			return
		case base + "socket":
			d.socket(w, r)
			return
		}
		respond(w, http.StatusNotFound)
	case http.MethodPost:
		d.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (d *dastVulnerabilityMessagesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	dastVulnerabilityMessages, err := dastVulnerabilityMessageRepository.find(makeQuery(r, dastVulnerabilityMessage{}, true), loadAllRelation)
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

func (d *dastVulnerabilityMessagesHandler) count(w http.ResponseWriter, r *http.Request) {
	dastVulnerabilityMessage, err := dastVulnerabilityMessageRepository.find(makeQuery(r, dastVulnerabilityMessage{}, true), loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondMessage(w, http.StatusOK, fmt.Sprintf("%d", len(dastVulnerabilityMessage)))
}

func (d *dastVulnerabilityMessagesHandler) fetchRange(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if emptyAny(start, end) {
		respond(w, http.StatusBadRequest)
		return
	}
	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	dastVulnerabilityMessage, err := dastVulnerabilityMessageRepository.findOrderLimit(makeQuery(r, dastVulnerabilityMessage{}, true), "created_at DESC", endInt, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if int64(len(dastVulnerabilityMessage)) < endInt {
		endInt = int64(len(dastVulnerabilityMessage))
	}
	response, err := reverse(dastVulnerabilityMessage[startInt:endInt])
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, response)
}

func (d *dastVulnerabilityMessagesHandler) socket(w http.ResponseWriter, r *http.Request) {
	putWebsocket(w, r, &dastVulnerabilityMessageWebsockets)
}

func (d *dastVulnerabilityMessagesHandler) store(w http.ResponseWriter, r *http.Request) {
	userID := r.PostFormValue("userId")
	vulnerabilityID := r.PostFormValue("vulnerabilityId")
	if emptyAny(userID, vulnerabilityID) {
		respond(w, http.StatusBadRequest)
		return
	}
	v, err := vulnerabilityRepository.first(map[string]interface{}{"id": vulnerabilityID}, vulnerabilityRelationScan, vulnerabilityRelationProjectUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, v.Project.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	dastVulnerabilityMessage, err := dastVulnerabilityMessageRepository.save(makeQuery(r, dastVulnerabilityMessage{}, false))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := dastVulnerabilityMessageRepository.findByID(dastVulnerabilityMessage.ID, loadAllRelation)
	project, err := projectRepository.findByID(response.Vulnerability.ProjectID, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	for _, u := range project.Users {
		socket, ok := dastVulnerabilityMessageWebsockets[u.ID]
		if !ok {
			continue
		}
		socket.WriteJSON(response)
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	go func() {
		u, err := authenticatedUser(r)
		if err != nil {
			return
		}
		place := "脆弱性(コミットのSHA1=" + v.Scan.CommitSHA1 + ")"
		sendMailMessage(u.Name, place, dastVulnerabilityMessage.Text, v.Project.Users)
	}()
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
	teamname := components[1]
	projectName := components[2]
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/"+path.Join(teamname, projectName, "info", "refs") {
			g.refs(w, r)
			return
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/" + path.Join(teamname, projectName, "git-receive-pack"):
			g.receivePack(w, r, teamname, projectName)
			return
		case "/" + path.Join(teamname, projectName, "git-upload-pack"):
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
	t, err := teamRepository.findByName(teamname)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if t.ID == 0 {
		respondError(w, http.StatusInternalServerError, errBadRequest)
		return
	}
	p, err := projectRepository.first(map[string]interface{}{"name": projectname, "team_id": t.ID})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if p.ID == 0 {
		respondError(w, http.StatusInternalServerError, errBadRequest)
		return
	}
	branchProtectionRules, err := branchProtectionRuleRepository.find(map[string]interface{}{"project_id": p.ID})
	if err != nil {
		return
	}
	branchName, commitSHA1, err := getBranchNameAndCommitSHA1(r)
	if err != nil {
		return
	}
	protection := false
	for _, branchProtectionRule := range branchProtectionRules {
		if branchName == branchProtectionRule.Branchname {
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
	meetings, err := meetingRepository.find(makeQuery(r, meeting{}, true), loadAllRelation)
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
	meetings, err := meetingRepository.findByIDs(stringsToUints(ids), loadAllRelation)
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
	projectID := r.PostFormValue("projectId")
	if emptyAny(projectID) {
		respond(w, http.StatusBadRequest)
		return
	}
	p, err := projectRepository.first(map[string]interface{}{"id": projectID}, projectRelationUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, p.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	meeting, err := meetingRepository.save(makeQuery(r, meeting{}, false))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := meetingRepository.findByID(meeting.ID, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, response)
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
		case base + "count":
			m.count(w, r)
			return
		case base + "range":
			m.fetchRange(w, r)
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
	meetingMessages, err := meetingMessageRepository.find(makeQuery(r, meetingMessage{}, true), loadAllRelation)
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

func (m *meetingMessagesHandler) count(w http.ResponseWriter, r *http.Request) {
	meetingMessages, err := meetingMessageRepository.find(makeQuery(r, meetingMessage{}, true), loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondMessage(w, http.StatusOK, fmt.Sprintf("%d", len(meetingMessages)))
}

func (m *meetingMessagesHandler) fetchRange(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if emptyAny(start, end) {
		respond(w, http.StatusBadRequest)
		return
	}
	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	meetingMessages, err := meetingMessageRepository.findOrderLimit(makeQuery(r, meetingMessage{}, true), "created_at DESC", endInt, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if int64(len(meetingMessages)) < endInt {
		endInt = int64(len(meetingMessages))
	}
	response, err := reverse(meetingMessages[startInt:endInt])
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, response)
}

func (m *meetingMessagesHandler) socket(w http.ResponseWriter, r *http.Request) {
	putWebsocket(w, r, &meetingMessageWebsockets)
}

func (m *meetingMessagesHandler) store(w http.ResponseWriter, r *http.Request) {
	meetingID := r.PostFormValue("meetingId")
	if emptyAny(meetingID) {
		respond(w, http.StatusBadRequest)
		return
	}
	meeting, err := meetingRepository.first(map[string]interface{}{"id": meetingID}, meetingRelationProjectUsers, meetingRelationUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, meeting.Project.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	meetingMessage, err := meetingMessageRepository.save(makeQuery(r, meetingMessage{}, false))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := meetingMessageRepository.findByID(meetingMessage.ID, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	for _, u := range meeting.Users {
		if u.ID == meetingMessage.UserID {
			continue
		}
		socket, ok := meetingMessageWebsockets[u.ID]
		if !ok {
			continue
		}
		socket.WriteJSON(response)
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	go func() {
		u, err := authenticatedUser(r)
		if err != nil {
			return
		}
		place := "ミーティング「" + meeting.Name + "」"
		sendMailMessage(u.Name, place, meetingMessage.Text, meeting.Users)
	}()
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
	meetingUsers, err := meetingUserRepository.find(makeQuery(r, meetingUser{}, true), loadAllRelation)
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
	meetingID := r.PostFormValue("meetingId")
	if emptyAny(meetingID) {
		respond(w, http.StatusBadRequest)
		return
	}
	meeting, err := meetingRepository.first(map[string]interface{}{"id": meetingID}, meetingRelationProjectUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, meeting.Project.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	meetingUser, err := meetingUserRepository.save(makeQuery(r, meetingUser{}, false))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	query := map[string]interface{}{"meeting_id": meetingUser.MeetingID, "user_id": meetingUser.UserID}
	response, err := meetingUserRepository.find(query, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, response)
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
	projects, err := projectRepository.find(makeQuery(r, project{}, true), loadAllRelation)
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
	ids := r.URL.Query()["ids"]
	if emptyAny(ids) {
		_, err := respondJSON(w, http.StatusOK, []project{})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	projects, err := projectRepository.findByIDs(stringsToUints(ids), loadAllRelation)
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
	teamID := r.PostFormValue("teamId")
	if emptyAny(teamID) {
		respond(w, http.StatusBadRequest)
		return
	}
	t, err := teamRepository.first(map[string]interface{}{"id": teamID}, teamRelationUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, t.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	project, err := projectRepository.save(makeQuery(r, project{}, false))
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
	response, err := projectRepository.findByID(project.ID, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, response)
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
	projectUsers, err := projectUserRepository.find(makeQuery(r, projectUser{}, true), loadAllRelation)
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
	project, err := projectRepository.first(map[string]interface{}{"id": projectID}, projectRelationTeamUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, project.Team.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	projectUserRole, err := projectUserRoleRepository.findByRole(role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	query := map[string]interface{}{"projectId": projectID, "roleId": projectUserRole.ID, "userId": userID}
	projectUser, err := projectUserRepository.save(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := projectUserRepository.first(map[string]interface{}{"project_id": projectUser.ProjectID, "user_id": projectUser.UserID})
	_, err = respondJSON(w, http.StatusInternalServerError, response)
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
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	handlename := r.PostFormValue("handlename")
	email := r.PostFormValue("email")
	if emptyAny(username, password, handlename, email) {
		respond(w, http.StatusNotFound)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	u, err := userRepository.saveWith(username, string(hashedPassword), handlename, email, pathNoImage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = login(w, username, password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, u)
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
	t, err := teamRepository.findByName(teamname)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := projectRepository.first(map[string]interface{}{"name": projectname, "team_id": t.ID})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	scans, err := scanRepository.find(map[string]interface{}{"project_id": p.ID}, loadAllRelation)
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
	hashedPassword := r.PostFormValue("password")
	teamname := r.PostFormValue("teamname")
	projectname := r.PostFormValue("projectname")
	if emptyAny(commitSHA1, username, hashedPassword, teamname, projectname) {
		respond(w, http.StatusBadRequest)
		return
	}
	u, err := userRepository.findByName(username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if u.Password != hashedPassword {
		respond(w, http.StatusBadRequest)
		return
	}
	t, err := teamRepository.findByName(teamname, teamRelationUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok := checkPermission(u, t.Users)
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	p, err := projectRepository.first(map[string]interface{}{"name": projectname, "team_id": t.ID})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	scan, err := scanRepository.saveWith(commitSHA1, p.ID, u.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := scanRepository.findByID(scan.ID, loadAllRelation)
	_, err = respondJSON(w, http.StatusOK, response)
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
	teams, err := teamRepository.find(makeQuery(r, team{}, true), loadAllRelation)
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
	teams, err := teamRepository.findByIDs(stringsToUints(ids), loadAllRelation)
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
	team, err := teamRepository.save(makeQuery(r, team{}, false))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := teamRepository.findByID(team.ID, loadAllRelation)
	_, err = respondJSON(w, http.StatusOK, response)
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
	teamUsers, err := teamUserRepository.find(makeQuery(r, teamUser{}, true), loadAllRelation)
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
	team, err := teamRepository.first(map[string]interface{}{"id": teamID}, teamRelationInvitationRequestsInviteeUser)
	u, err := authenticatedUser(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if team.FounderUserID != u.ID {
		users := make([]user, len(team.InvitationRequests))
		for index, invitationRequest := range team.InvitationRequests {
			users[index] = invitationRequest.InviteeUser
		}
		ok, err := checkPermissionWithRequest(r, users)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if !ok {
			respond(w, http.StatusForbidden)
			return
		}
	}
	teamUserRole, err := teamUserRoleRepository.findByRole(role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	query := map[string]interface{}{"teamID": teamID, "roleID": teamUserRole.ID, "userID": userID}
	teamUser, err := teamUserRepository.save(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := teamUserRepository.first(map[string]interface{}{"team_id": teamUser.TeamID, "user_id": teamUser.UserID})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusInternalServerError, response)
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
	teamUserInvitationRequestProjects, err := teamUserInvitationRequestProjectRepository.find(makeQuery(r, teamUserInvitationRequestProject{}, true), loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, teamUserInvitationRequestProjects)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUserInvitationRequestProjectsHandler) store(w http.ResponseWriter, r *http.Request) {
	teamUserInvitationRequestID := r.PostFormValue("teamUserInvitationRequestId")
	if emptyAny(teamUserInvitationRequestID) {
		respond(w, http.StatusBadRequest)
		return
	}
	teamUserInvitationRequest, err := teamUserInvitationRequestRepository.first(map[string]interface{}{"id": teamUserInvitationRequestID}, teamUserInvitationRequestRelationTeamUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, teamUserInvitationRequest.Team.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	teamUserInvitationRequestProject, err := teamUserInvitationRequestProjectRepository.save(makeQuery(r, teamUserInvitationRequestProject{}, false))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	query := map[string]interface{}{"team_user_invitation_request_id": teamUserInvitationRequestProject.TeamUserInvitationRequestID, "project_id": teamUserInvitationRequestProject.ProjectID}
	response, err := teamUserInvitationRequestProjectRepository.first(query, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUserInvitationRequestProjectsHandler) destroy(w http.ResponseWriter, r *http.Request) {
	err := teamUserInvitationRequestProjectRepository.delete(makeQuery(r, teamUserInvitationRequestProject{}, true))
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
	teamUserInvitationRequests, err := teamUserInvitationRequestRepository.find(makeQuery(r, teamUserInvitationRequest{}, true), loadAllRelation)
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
	if emptyAny(inviterUserID, inviteeUserID, roleText, teamID) {
		respond(w, http.StatusBadRequest)
		return
	}
	team, err := teamRepository.first(map[string]interface{}{"id": teamID}, teamRelationUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, team.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	role, err := teamUserRoleRepository.findByRole(roleText)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	query := map[string]interface{}{"message": message, "inviterUserID": inviterUserID, "inviteeUserID": inviteeUserID, "roleID": role.ID, "teamID": teamID}
	teamUserInvitationRequest, err := teamUserInvitationRequestRepository.save(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := teamUserInvitationRequestRepository.findByID(teamUserInvitationRequest.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *teamUserInvitationRequestsHandler) destroy(w http.ResponseWriter, r *http.Request) {
	err := teamUserInvitationRequestRepository.delete(makeQuery(r, teamUserInvitationRequest{}, true))
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
	tests, err := testRepository.find(makeQuery(r, test{}, true), loadAllRelation)
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
	putWebsocket(w, r, &testWebsockets)
}

func (t *testsHandler) revision(w http.ResponseWriter, r *http.Request) {
	teamname := r.URL.Query().Get("teamname")
	projectname := r.URL.Query().Get("projectname")
	revision := r.URL.Query().Get("revision")
	if emptyAny(teamname, projectname, revision) {
		respond(w, http.StatusBadRequest)
		return
	}
	team, err := teamRepository.findByName(teamname)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := projectRepository.first(map[string]interface{}{"name": projectname, "team_id": team.ID})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	tests, err := testRepository.findOrder(map[string]interface{}{"project_id": p.ID}, "created_at DESC", loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
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
		case base + "count":
			t.count(w, r)
			return
		case base + "range":
			t.fetchRange(w, r)
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
	testMessages, err := testMessageRepository.find(makeQuery(r, testMessage{}, true), loadAllRelation)
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

func (t *testMessagesHandler) count(w http.ResponseWriter, r *http.Request) {
	testMessages, err := testMessageRepository.find(makeQuery(r, testMessage{}, true), loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondMessage(w, http.StatusOK, fmt.Sprintf("%d", len(testMessages)))
}

func (t *testMessagesHandler) fetchRange(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if emptyAny(start, end) {
		respond(w, http.StatusBadRequest)
		return
	}
	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	testMessages, err := testMessageRepository.findOrderLimit(makeQuery(r, testMessage{}, true), "created_at DESC", endInt, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if int64(len(testMessages)) < endInt {
		endInt = int64(len(testMessages))
	}
	response, err := reverse(testMessages[startInt:endInt])
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, response)
}

func (t *testMessagesHandler) socket(w http.ResponseWriter, r *http.Request) {
	putWebsocket(w, r, &testMessageWebsockets)
}

func (t *testMessagesHandler) store(w http.ResponseWriter, r *http.Request) {
	testID := r.PostFormValue("testId")
	if emptyAny(testID) {
		respond(w, http.StatusBadRequest)
		return
	}
	test, err := testRepository.first(map[string]interface{}{"id": testID}, testRelationProjectUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := checkPermissionWithRequest(r, test.Project.Users)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	testMessage, err := testMessageRepository.save(makeQuery(r, testMessage{}, false))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := testMessageRepository.findByID(testMessage.ID, loadAllRelation)
	for _, u := range test.Project.Users {
		socket, ok := testMessageWebsockets[u.ID]
		if !ok {
			continue
		}
		socket.WriteJSON(response)
	}
	_, err = respondJSON(w, http.StatusOK, response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	go func() {
		u, err := authenticatedUser(r)
		if err != nil {
			return
		}
		place := "テスト(コミットのSHA1=" + test.CommitSHA1 + ")"
		sendMailMessage(u.Name, place, testMessage.Text, test.Project.Users)
	}()
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
	testResults, err := testResultRepository.find(makeQuery(r, testResult{}, true), loadAllRelation)
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
	testStatuses, err := testStatusRepository.find(makeQuery(r, testStatus{}, true))
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
	id, err := session.Uint(sessionStoreNameUserID)
	user, err := userRepository.findByID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
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
	users, err := userRepository.find(makeQuery(r, user{}, true))
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
	users, err := userRepository.findByIDs(stringsToUints(ids))
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
	vulnerabilities, err := vulnerabilityRepository.find(makeQuery(r, vulnerability{}, true), loadAllRelation)
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
	username := r.PostFormValue("username")
	hashedPassword := r.PostFormValue("password")
	if emptyAny(name, description, path, method, request, response, scanID, teamname, projectname, username, hashedPassword) {
		respond(w, http.StatusBadRequest)
		return
	}
	u, err := userRepository.findByName(username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if u.Password != hashedPassword {
		respond(w, http.StatusBadRequest)
		return
	}
	t, err := teamRepository.findByName(teamname, teamRelationUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	ok := checkPermission(u, t.Users)
	if !ok {
		respond(w, http.StatusForbidden)
		return
	}
	p, err := projectRepository.first(map[string]interface{}{"name": projectname, "team_id": t.ID})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	query := map[string]interface{}{"Name": name, "Description": description, "Path": path, "Method": method, "Request": request, "Response": response, "ProjectID": p.ID, "ScanID": scanID}
	vulnerability, err := vulnerabilityRepository.save(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	res, err := vulnerabilityRepository.findByID(vulnerability.ID, loadAllRelation)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = respondJSON(w, http.StatusOK, res)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
}
