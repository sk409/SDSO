package main

import (
	"net/http"

	"github.com/sk409/goconst"
)

func main() {
	http.Handle("/", gitBasicAuth(&gitHandler{}))
	http.Handle("/auth", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodGet}, preflight(&authHandler{}))))))
	http.Handle("/branches", cors(&branchesHandler{}))
	http.Handle("/branch_protection_rules", cors(&branchProtectionRulesHandler{}))
	http.Handle("/commits/", cors(&commitsHandler{}))
	http.Handle("/files/", cors(&filesHandler{}))
	http.Handle("/login", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodPost}, preflight(&loginHandler{}))))))
	http.Handle("/logout", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodPost}, preflight(&logoutHandler{}))))))
	http.Handle("/projects", cors(&projectsHandler{}))
	http.Handle("/register", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodPost}, preflight(&registerHandler{}))))))
	http.Handle("/repositories", cors(&repositoriesHandler{}))
	http.Handle("/scans", cors(&scansHandler{}))
	http.Handle("/tests/", cors(&testsHandler{}))
	http.Handle("/test_results", cors(&testResultsHandler{}))
	http.Handle("/test_statuses", cors(&testStatusesHandler{}))
	http.Handle("/user", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodGet}, preflight(&userHandler{}))))))
	http.Handle("/users", cors(&usersHandler{}))
	http.Handle("/vulnerabilities", cors(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE}, allowMethods([]string{http.MethodGet, http.MethodPost}, preflight(&vulnerabilitiesHandler{})))))

	// authRouter := mux.NewRouter().PathPrefix("/auth")
	// authCheckRouter := authRouter.Subrouter()
	// authCheckRouter.Use(cors)
	// authCheckRouter.HandleFunc("/check", authCheckHandler).Methods(http.MethodGet, http.MethodOptions)
	// http.Handle("/auth/check", authCheckRouter)

	// loginRouter := mux.NewRouter()
	// loginRouter.Use(cors)
	// loginRouter.HandleFunc("/login", loginHandler).Methods(http.MethodPost, http.MethodOptions)
	// http.Handle("/login", loginRouter)

	// socialLoginRouter := mux.NewRouter()
	// socialLoginRouter.Use(cors)
	// socialLoginRouter.HandleFunc("/social_login", socialLoginHandler).Methods(http.MethodGet)
	// http.Handle("/social_login", socialLoginRouter)

	// /*************************************/
	// // socialLoginCallbackRouter := mux.NewRouter()
	// // socialLoginCallbackRouter.Use(cors)
	// // socialLoginCallbackRouter.HandleFunc("/auth/callback/{provider}", socialLoginCallbackHandler)
	// // http.Handle("/auth/callback/{provider}", socialLoginCallbackRouter)
	// http.HandleFunc("/auth/callback/", socialLoginCallbackHandler)
	// /*************************************/

	// logoutRouter := mux.NewRouter()
	// logoutRouter.Use(cors)
	// logoutRouter.HandleFunc("/logout", logoutHandler).Methods(http.MethodPost, http.MethodOptions)
	// http.Handle("/logout", logoutRouter)

	// registerRouter := mux.NewRouter()
	// registerRouter.Use(cors)
	// registerRouter.HandleFunc("/register", registerHandler).Methods(http.MethodPost, http.MethodOptions)
	// http.Handle("/register", registerRouter)

	// userRouter := mux.NewRouter()
	// userRouter.Use(cors)
	// userRouter.HandleFunc("/user", fetchUserHandler).Methods(http.MethodGet, http.MethodOptions)
	// http.Handle("/user", userRouter)

	// usersRouter := mux.NewRouter()
	// usersRouter.Use(cors)
	// usersRouter.HandleFunc("/users", fetchUsersHandler).Methods(http.MethodGet)
	// http.Handle("/users", usersRouter)

	// usersExistRouter := mux.NewRouter()
	// usersExistRouter.Use(cors)
	// usersExistRouter.HandleFunc("/users/exist", existUserHandler).Methods(http.MethodGet)
	// http.Handle("/users/exist", usersExistRouter)

	// projectsRouter := mux.NewRouter()
	// projectsRouter.Use(cors)
	// projectsRouter.HandleFunc("/projects", fetchProjectsHandler).Methods(http.MethodGet)
	// projectsRouter.HandleFunc("/projects", storeProjectHandler).Methods(http.MethodPost, http.MethodOptions)
	// http.Handle("/projects", projectsRouter)

	// projectsExistHandler := mux.NewRouter()
	// projectsExistHandler.Use(cors)
	// projectsExistHandler.HandleFunc("/projects/exist", existProjectHandler).Methods(http.MethodGet)
	// http.Handle("/projects/exist", projectsExistHandler)

	// vulnerabilitiesRouter := mux.NewRouter()
	// vulnerabilitiesRouter.Use(cors)
	// vulnerabilitiesRouter.HandleFunc("/vulnerabilities", fetchVulnerabilities).Methods(http.MethodGet)
	// vulnerabilitiesRouter.HandleFunc("/vulnerabilities", storeVulnerability).Methods(http.MethodPost, http.MethodOptions)
	// http.Handle("/vulnerabilities", vulnerabilitiesRouter)

	// scansRouter := mux.NewRouter()
	// scansRouter.Use(cors)
	// scansRouter.HandleFunc("/scans", fetchScansHandler).Methods(http.MethodGet)
	// scansRouter.HandleFunc("/scans", storeScanHandler).Methods(http.MethodPost, http.MethodOptions)
	// http.Handle("/scans", scansRouter)

	// repositoriesRouter := mux.NewRouter()
	// repositoriesRouter.Use(cors)
	// repositoriesRouter.HandleFunc("/repositories/init", initRepositoryHandler).Methods(http.MethodPost, http.MethodOptions)
	// http.Handle("/repositories/", repositoriesRouter)

	// filesRouter := mux.NewRouter()
	// filesRouter.Use(cors)
	// filesRouter.HandleFunc("/file", fetchFilesHandler).Methods(http.MethodGet)
	// // filesRouter.HandleFunc("/files/text", fetchFileTextHandler).Methods(http.MethodGet)
	// http.Handle("/file", filesRouter)

	// filesTextRouter := mux.NewRouter()
	// filesTextRouter.Use(cors)
	// filesTextRouter.HandleFunc("/files/text", fetchFileTextHandler).Methods(http.MethodGet)
	// http.Handle("/files/text", filesTextRouter)

	// testsRouter := mux.NewRouter()
	// testsRouter.Use(cors)
	// testsRouter.HandleFunc("/tests", fetchTestsHandler).Methods(http.MethodGet)
	// http.Handle("/tests", testsRouter)

	// testingRouter := mux.NewRouter()
	// testingRouter.Use(cors)
	// testingRouter.Schemes("ws")
	//testingRouter.HandleFunc("/testing", testingHandler)
	http.HandleFunc("/test_socket", testSocketHandler)
	http.HandleFunc("/test_result_socket", testResultSocketHandler)

	// testResultsRouter := mux.NewRouter()
	// testResultsRouter.Use(cors)
	// testResultsRouter.HandleFunc("/test_results", fetchTestResultsHandler).Methods(http.MethodGet)
	// http.Handle("/test_results", testResultsRouter)

	// testStatusesRouter := mux.NewRouter()
	// testStatusesRouter.Use(cors)
	// testStatusesRouter.HandleFunc("/test_statuses", fetchTestStatuses).Methods(http.MethodGet)
	// http.Handle("/test_statuses", testStatusesRouter)

	// branchProtectionRulesRouter := mux.NewRouter()
	// branchProtectionRulesRouter.Use(cors)
	// branchProtectionRulesRouter.HandleFunc("/branch_protection_rules", fetchBranchProtectionRules).Methods(http.MethodGet)
	// branchProtectionRulesRouter.HandleFunc("/branch_protection_rules", storeBranchProtectionRules).Methods(http.MethodPost)
	// http.Handle("/branch_protection_rules", branchProtectionRulesRouter)

	// branchesRouter := mux.NewRouter()
	// branchesRouter.Use(cors)
	// branchesRouter.HandleFunc("/branches", fetchBranchesHandler).Methods(http.MethodGet)
	// http.Handle("/branches", branchesRouter)

	// commitSHA1sRouter := mux.NewRouter()
	// commitSHA1sRouter.Use(cors)
	// commitSHA1sRouter.HandleFunc("/commit_sha1s", fetchCommitSHA1sHandler).Methods(http.MethodGet)
	// http.Handle("/commit_sha1s", commitSHA1sRouter)

	// commitsRouter := mux.NewRouter()
	// commitsRouter.Use(cors)
	// commitsRouter.HandleFunc("/commits", fetchCommitsHandler).Methods(http.MethodGet)
	// http.Handle("/commits", commitsRouter)

	// commitsShowRouter := mux.NewRouter()
	// commitsShowRouter.Use(cors)
	// commitsShowRouter.HandleFunc("/commits/show", showCommitHandler).Methods(http.MethodGet)
	// http.Handle("/commits/show", commitsShowRouter)

	//http.Handle("/commits/", cors(commitsHandler{}))

	// gitRouter := mux.NewRouter()
	// gitRouter.HandleFunc("/{user}/{project}/info/refs", gitInfoRefsHandler).Methods(http.MethodGet)
	// gitRouter.Handle("/{user}/{project}/git-receive-pack", basicAuth(http.HandlerFunc(gitReceivePackHandler))).Methods(http.MethodPost)
	// //gitRouter.HandleFunc("/{user}/{project}/git-receive-pack", gitReceivePackHandler).Methods(http.MethodPost)
	// gitRouter.HandleFunc("/{user}/{project}/git-upload-pack", gitUploadPackHandler).Methods(http.MethodPost)
	// http.Handle("/", gitRouter)

	http.ListenAndServe(serverHostAndPort, nil)
}
