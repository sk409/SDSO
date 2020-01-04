package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"

	"github.com/gorilla/mux"

	"github.com/sk409/goconst"
)

func makeQueryAndValues(r *http.Request) (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}
	if r.Method == http.MethodGet {
		for key, value := range r.URL.Query() {
			conditions = append(conditions, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	} else {
		r.ParseForm()
		for key, value := range r.PostForm {
			conditions = append(conditions, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	}
	query := strings.Join(conditions, " AND ")
	return query, values
}

func existHandler(w http.ResponseWriter, r *http.Request, tableName string) {
	query, values := makeQueryAndValues(r)
	count := 0
	db.Table(tableName).Where(query, values...).Count(&count)
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonString := fmt.Sprintf("{\"exist\":%t}", count != 0)
	w.Write([]byte(jsonString))
}

func socketHandler(w http.ResponseWriter, r *http.Request, websockets *map[uint]*websocket.Conn) {
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionID := sessionCookie.Value
	session, err := sessionManager.Provider.Get(sessionID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := session.Uint(sessionStoreNameUserID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	socket, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	(*websockets)[userID] = socket
}

func authCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("authCheckHandler")
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodGet)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	notAuthenticated := []byte("{\"authenticated\": false}")
	authenticated := []byte("{\"authenticated\": true}")
	sessionCookie, err := r.Cookie("sessionID")
	if err != nil {
		w.Write(notAuthenticated)
		return
	}
	session, err := sessionManager.Provider.Get(sessionCookie.Value)
	if session == nil || err != nil {
		w.Write(notAuthenticated)
		http.SetCookie(w, &http.Cookie{
			Name:   cookieNameSessionID,
			MaxAge: -1,
		})
		return
	}
	w.Write(authenticated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginHandler")
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodPost)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	if r.Method == http.MethodOptions {
		return
	}
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")
	if len(name) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(password)))
	user := user{}
	db.Model(user).Where("name = ? AND password = ?", name, hashedPassword).First(&user)
	exist := user.ID != 0
	if exist {
		session, err := sessionManager.Provider.Start()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		session.Store(sessionStoreNameUserID, user.ID)
		sessionCookie := newCookie(cookieNameSessionID, session.ID(), cookie30Days)
		http.SetCookie(w, sessionCookie)
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	response := fmt.Sprintf("{\"exist\": %t}", exist)
	w.Write([]byte(response))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logoutHandler")
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodGet)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	if r.Method == http.MethodOptions {
		return
	}
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err == nil {
		sessionManager.Provider.Stop(sessionCookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   cookieNameSessionID,
		MaxAge: -1,
	})
}

func socialLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("socialLoginHandler")
	providerName := r.URL.Query().Get("provider")
	if len(providerName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginURL, err := provider.GetBeginAuthURL(nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_LOCATION, loginURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func socialLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	providerName := segs[3]
	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userInfo, err := provider.GetUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email := userInfo.Email()
	usr := user{}
	db.Where("name = ?", email).First(&usr)
	if usr.ID == 0 {
		profileImagePath := userInfo.AvatarURL()
		usr = user{Name: email, ProfileImagePath: &profileImagePath}
		db.Save(&usr)
	}
	session, err := sessionManager.Provider.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Store(sessionStoreNameUserID, usr.ID)
	sessionCookie := newCookie(cookieNameSessionID, session.ID(), cookie30Days)
	http.SetCookie(w, sessionCookie)
	w.Header().Set(goconst.HTTP_HEADER_LOCATION, "http://localhost:3000")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("registerHandler")
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodPost)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	if r.Method == http.MethodOptions {
		return
	}
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")
	if len(name) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(password)))
	user := user{Name: name, Password: hashedPassword}
	count := 0
	db.Model(user).Where("name = ? AND password = ?", user.Name, user.Password).Count(&count)
	if count != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	db.Save(&user)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session, err := sessionManager.Provider.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Store(sessionStoreNameUserID, user.ID)
	sessionCookie := newCookie(cookieNameSessionID, session.ID(), cookie30Days)
	http.SetCookie(w, sessionCookie)
}

func fetchUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchUserHandler")
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodGet)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	if r.Method == http.MethodOptions {
		return
	}
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session, err := sessionManager.Provider.Get(sessionCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := session.Uint(sessionStoreNameUserID)
	user := user{}
	user.ID = userID
	db.Find(&user)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func fetchUsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchUsersHandler")
	query, values := makeQueryAndValues(r)
	users := []user{}
	db.Where(query, values...).Find(&users)
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func existUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("existUserHandler")
	existHandler(w, r, tableNameUsers)
}

func storeProjectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("storeProjectHandler")
	r.ParseForm()
	// log.Println(r.PostForm)
	name := r.PostFormValue("name")
	userIDString := r.PostFormValue("userID")
	if len(name) == 0 || len(userIDString) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.Save(&project{Name: name, UserID: uint(userID)})
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func fetchProjectsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchProjectsHandler")
	query, values := makeQueryAndValues(r)
	projects := []project{}
	db.Model(project{}).Where(query, values...).Find(&projects)
	jsonBytes, err := json.Marshal(projects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func existProjectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("existProjectHandler")
	existHandler(w, r, tableNameProjects)
}

func storeScanHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("storeScanHandler")
	if r.Method == http.MethodOptions {
		return
	}
	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")
	projectName := r.PostFormValue("projectName")
	if len(userName) == 0 || len(password) == 0 || len(projectName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := user{}
	db.Where("name = ? AND password = ?", userName, password).First(&user)
	if db.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userIDString := strconv.Itoa(int(user.ID))
	project := project{}
	db.Where("name = ? AND user_id = ?", projectName, userIDString).First(&project)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	scan := scan{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	db.Save(&scan)
	jsonBytes, err := json.Marshal(scan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func fetchScansHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchScansHandler")
	query, values := makeQueryAndValues(r)
	scans := []scan{}
	db.Where(query, values...).Find(&scans)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonBytes, err := json.Marshal(scans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func storeVulnerability(w http.ResponseWriter, r *http.Request) {
	fmt.Println("storeVulnerability")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodPost)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, goconst.HTTP_HEADER_CONTENT_TYPE)
	if r.Method == http.MethodOptions {
		return
	}
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
	if len(name) == 0 ||
		len(description) == 0 ||
		len(path) == 0 ||
		len(method) == 0 ||
		len(request) == 0 ||
		len(response) == 0 ||
		len(scanID) == 0 ||
		len(userName) == 0 ||
		len(password) == 0 ||
		len(projectName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	scanIDInt, err := strconv.Atoi(scanID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := &user{}
	db.Model(user).Where("name = ? AND password = ?", userName, password).First(&user)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	project := &project{}
	db.Model(project).Where("name = ? AND user_id = ?", projectName, user.ID).First(&project)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vulnerability := vulnerability{
		Name:        name,
		Description: description,
		Path:        path,
		Method:      method,
		Request:     request,
		Response:    response,
		ProjectID:   project.ID,
		ScanID:      uint(scanIDInt),
	}
	db.Save(&vulnerability)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func fetchVulnerabilities(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchVulnerabilities")
	query, values := makeQueryAndValues(r)
	vulnerabilities := []vulnerability{}
	db.Where(query, values...).Find(&vulnerabilities)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonBytes, err := json.Marshal(vulnerabilities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func initRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("initRepositoryHandler")
	userName := r.PostFormValue("userName")
	projectName := r.PostFormValue("projectName")
	if len(userName) == 0 || len(projectName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	repositoryPath := filepath.Join(userName, projectName)
	exist := existDirectory(filepath.Join(gitRepositories.RootDirectoryPath, repositoryPath))
	// if err != nil {
	// 	log.Println(err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	if exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//mkDirIfNotExist(filepath.Join(gitRepositories.RootPath, repositoryPath))
	err := os.MkdirAll(filepath.Join(gitRepositories.RootDirectoryPath, repositoryPath), 0755)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = gitRepositories.InitBare(repositoryPath)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func fetchFileTextHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchFileTextHandler")
	userName := r.URL.Query().Get("userName")
	projectName := r.URL.Query().Get("projectName")
	path := r.URL.Query().Get("path")
	if len(userName) == 0 || len(projectName) == 0 || len(path) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, err := os.Open(filepath.Join(gitClones.RootDirectoryPath, userName, projectName, path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileTextBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_PLAIN_TEXT)
	w.Write(fileTextBytes)
}

func fetchFilesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchFilesHandler")
	userName := r.URL.Query().Get("userName")
	projectName := r.URL.Query().Get("projectName")
	path := r.URL.Query().Get("path")
	output, err := gitClones.LsFiles(filepath.Join(userName, projectName))
	if err != nil {
		log.Println(err.Error())
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}
	outputLines := strings.Split(string(output), "\n")
	files := make(map[string]map[string]interface{})
	for _, outputLine := range outputLines {
		var fileName string
		if len(path) == 0 {
			// TODO: OSごとにデリミタを変える?
			fileName = strings.Split(outputLine, "/")[0]
		} else if strings.HasPrefix(outputLine, path) {
			// TODO: 範囲外アクセスにならない?
			fileName = strings.Split(strings.TrimPrefix(outputLine, path), "/")[1]
		}
		if len(fileName) == 0 {
			continue
		}
		_, exist := files[fileName]
		if exist {
			continue
		}
		directory, err := isDirectory(filepath.Join(gitClones.RootDirectoryPath, userName, projectName, path, fileName))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		file := map[string]interface{}{}
		file["path"] = filepath.Join(path, fileName)
		file["name"] = fileName
		file["isDirectory"] = directory
		files[fileName] = file
	}
	jsonBytes, err := json.Marshal(&files)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func fetchTestsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchTestsHandler")
	query, values := makeQueryAndValues(r)
	tests := []test{}
	db.Where(query, values...).Find(&tests)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(tests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func testSocketHandler(w http.ResponseWriter, r *http.Request) {
	socketHandler(w, r, &websocketsTest)
}

func testResultSocketHandler(w http.ResponseWriter, r *http.Request) {
	socketHandler(w, r, &websocketsTestResult)
	// c := make(chan []byte)
	// read := func() {
	// 	for {
	// 		if _, msg, err := socket.ReadMessage(); err == nil {
	// 			c <- msg
	// 		} else {
	// 			break
	// 		}
	// 	}
	// 	socket.Close()
	// }
	// write := func() {
	// 	for msg := range c {
	// 		if err := socket.WriteMessage(websocket.TextMessage, msg); err != nil {
	// 			break
	// 		}
	// 	}
	// 	socket.Close()
	// }
	// go write()
	// read()
}

func fetchTestResultsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchTestResultsHandler")
	query, values := makeQueryAndValues(r)
	testResults := []testResult{}
	db.Where(query, values...).Find(&testResults)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(testResults)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func fetchTestStatuses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchTestStatuses")
	query, values := makeQueryAndValues(r)
	testStatuses := []testStatus{}
	db.Where(query, values...).Find(&testStatuses)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(testStatuses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func gitInfoRefsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("gitInfoRefsHandler")
	gitServer.ServeHTTP(w, r)
}

func gitReceivePackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("gitReceivePackHandler")
	pathParams := mux.Vars(r)
	userName := pathParams["user"]
	projectName := pathParams["project"]
	gitServer.ServeHTTP(w, r)
	repositoryPath := serverScheme + "://" + path.Join(serverHostAndPort, userName, projectName)
	clonePath := filepath.Join(gitClones.RootDirectoryPath, filepath.Join(userName, projectName))
	if existDirectory(clonePath) {
		os.RemoveAll(clonePath)
	}
	gitClones.Clone(
		repositoryPath,
		filepath.Join(userName, projectName),
	)
	go runTest(userName, projectName)
}

func gitUploadPackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("gitUploadPackHandler")
	gitServer.ServeHTTP(w, r)
}
