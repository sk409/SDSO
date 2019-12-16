package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func f(w http.ResponseWriter, r *http.Request) {
	// log.Println("OK")
	// log.Println(r.Method)
	// // w.Header().Set("Access-Control-Allow-Origin", "*")
	// log.Println(r.RemoteAddr)
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func main() {

	authRouter := mux.NewRouter().PathPrefix("/auth")
	authCheckRouter := authRouter.Subrouter()
	authCheckRouter.Use(corsMiddleware)
	authCheckRouter.HandleFunc("/check", authCheckHandler).Methods(http.MethodGet, http.MethodOptions)
	http.Handle("/auth/check", authCheckRouter)

	loginRouter := mux.NewRouter()
	loginRouter.Use(corsMiddleware)
	loginRouter.HandleFunc("/login", loginHandler).Methods(http.MethodPost, http.MethodOptions)
	http.Handle("/login", loginRouter)

	registerRouter := mux.NewRouter()
	registerRouter.Use(corsMiddleware)
	registerRouter.HandleFunc("/register", registerHandler).Methods(http.MethodPost, http.MethodOptions)
	http.Handle("/register", registerRouter)

	userRouter := mux.NewRouter()
	userRouter.Use(corsMiddleware)
	userRouter.HandleFunc("/user", fetchUserHandler).Methods(http.MethodGet, http.MethodOptions)
	http.Handle("/user", userRouter)

	usersRouter := mux.NewRouter()
	usersRouter.Use(corsMiddleware)
	usersRouter.HandleFunc("/users", fetchUsersHandler).Methods(http.MethodGet)
	http.Handle("/users", usersRouter)

	usersExistRouter := mux.NewRouter()
	usersExistRouter.Use(corsMiddleware)
	usersExistRouter.HandleFunc("/users/exist", existUserHandler).Methods(http.MethodGet)
	http.Handle("/users/exist", usersExistRouter)

	projectsRouter := mux.NewRouter()
	projectsRouter.Use(corsMiddleware)
	projectsRouter.HandleFunc("/projects", fetchProjectsHandler).Methods(http.MethodGet)
	projectsRouter.HandleFunc("/projects", storeProjectHandler).Methods(http.MethodPost, http.MethodOptions)
	http.Handle("/projects", projectsRouter)

	projectsExistHandler := mux.NewRouter()
	projectsExistHandler.Use(corsMiddleware)
	projectsExistHandler.HandleFunc("/projects/exist", existProjectHandler).Methods(http.MethodGet)
	http.Handle("/projects/exist", projectsExistHandler)

	vulnerabilitiesRouter := mux.NewRouter()
	vulnerabilitiesRouter.Use(corsMiddleware)
	vulnerabilitiesRouter.HandleFunc("/vulnerabilities", fetchVulnerabilities).Methods(http.MethodGet)
	vulnerabilitiesRouter.HandleFunc("/vulnerabilities", storeVulnerability).Methods(http.MethodPost, http.MethodOptions)
	http.Handle("/vulnerabilities", vulnerabilitiesRouter)

	scansRouter := mux.NewRouter()
	scansRouter.Use(corsMiddleware)
	scansRouter.HandleFunc("/scans", fetchScansHandler).Methods(http.MethodGet)
	scansRouter.HandleFunc("/scans", storeScanHandler).Methods(http.MethodPost, http.MethodOptions)
	http.Handle("/scans", scansRouter)

	repositoriesRouter := mux.NewRouter()
	repositoriesRouter.Use(corsMiddleware)
	repositoriesRouter.HandleFunc("/repositories/init", initRepositoryHandler).Methods(http.MethodPost, http.MethodOptions)
	http.Handle("/repositories/", repositoriesRouter)

	filesRouter := mux.NewRouter()
	filesRouter.Use(corsMiddleware)
	filesRouter.HandleFunc("/file", fetchFilesHandler).Methods(http.MethodGet)
	// filesRouter.HandleFunc("/files/text", fetchFileTextHandler).Methods(http.MethodGet)
	http.Handle("/file", filesRouter)

	filesTextRouter := mux.NewRouter()
	filesTextRouter.Use(corsMiddleware)
	filesTextRouter.HandleFunc("/files/text", fetchFileTextHandler).Methods(http.MethodGet)
	http.Handle("/files/text", filesTextRouter)

	testsRouter := mux.NewRouter()
	testsRouter.Use(corsMiddleware)
	testsRouter.HandleFunc("/tests", fetchTestsHandler).Methods(http.MethodGet)
	http.Handle("/tests", testsRouter)

	testResultsRouter := mux.NewRouter()
	testResultsRouter.Use(corsMiddleware)
	testResultsRouter.HandleFunc("/test_results", fetchTestResultsHandler).Methods(http.MethodGet)
	http.Handle("/test_results", testResultsRouter)

	testStatusesRouter := mux.NewRouter()
	testStatusesRouter.Use(corsMiddleware)
	testStatusesRouter.HandleFunc("/test_statuses", fetchTestStatuses).Methods(http.MethodGet)
	http.Handle("/test_statuses", testStatusesRouter)

	gitRouter := mux.NewRouter()
	gitRouter.HandleFunc("/{user}/{project}/info/refs", gitInfoRefsHandler).Methods(http.MethodGet)
	gitRouter.HandleFunc("/{user}/{project}/git-receive-pack", gitReceivePackHandler).Methods(http.MethodPost)
	gitRouter.HandleFunc("/{user}/{project}/git-upload-pack", gitUploadPackHandler).Methods(http.MethodPost)
	http.Handle("/", gitRouter)

	http.ListenAndServe(serverHostAndPort, nil)
}
