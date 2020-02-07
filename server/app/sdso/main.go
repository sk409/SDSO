package main

import (
	"net/http"

	"github.com/sk409/goconst"
)

func main() {
	// http.Handle("/", gitBasicAuth(&gitHandler{}))
	rootRouter := router{}
	rootRouter.gitBasicAuth()
	rootRouter.handler = &gitHandler{}
	http.Handle("/", &rootRouter)
	// http.Handle("/auth", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodGet}, preflight(&authHandler{}))))))
	authRouter := router{}
	authRouter.cors()
	authRouter.allowCredentials(http.MethodGet)
	authRouter.allowHeaders(map[string][]string{http.MethodGet: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	authRouter.allowMethods(http.MethodGet)
	authRouter.preflight()
	authRouter.handler = &authHandler{}
	http.Handle("/auth", &authRouter)
	// http.Handle("/branches", cors(&branchesHandler{}))
	branchesRouter := router{}
	branchesRouter.cors()
	branchesRouter.handler = &branchesHandler{}
	http.Handle("/branches", &branchesRouter)
	// http.Handle("/branch_protection_rules", cors(&branchProtectionRulesHandler{}))
	branchProtectionRulesRouter := router{}
	branchProtectionRulesRouter.cors()
	branchProtectionRulesRouter.handler = &branchProtectionRulesHandler{}
	http.Handle("/branch_protection_rules", &branchProtectionRulesRouter)
	// http.Handle("/commits/", cors(&commitsHandler{}))
	commitsRouter := router{}
	commitsRouter.cors()
	commitsRouter.handler = &commitsHandler{}
	http.Handle("/commits/", &commitsRouter)
	//
	dastVulnerabilityMessagesRouter := router{}
	dastVulnerabilityMessagesRouter.cors()
	dastVulnerabilityMessagesRouter.handler = &dastVulnerabilityMessagesHandler{}
	http.Handle("/dast_vulnerability_messages/", &dastVulnerabilityMessagesRouter)
	// http.Handle("/files/", cors(&filesHandler{}))
	filesRouter := router{}
	filesRouter.cors()
	filesRouter.handler = &filesHandler{}
	http.Handle("/files/", &filesRouter)
	// http.Handle("/login", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodPost}, preflight(&loginHandler{}))))))
	loginRouter := router{}
	loginRouter.cors()
	loginRouter.allowCredentials(http.MethodPost)
	loginRouter.allowHeaders(map[string][]string{http.MethodPost: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	loginRouter.allowMethods(http.MethodPost)
	loginRouter.preflight()
	loginRouter.handler = &loginHandler{}
	http.Handle("/login", &loginRouter)
	// http.Handle("/logout", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodPost}, preflight(&logoutHandler{}))))))
	logoutRouter := router{}
	logoutRouter.cors()
	logoutRouter.allowCredentials(http.MethodPost)
	logoutRouter.allowHeaders(map[string][]string{http.MethodPost: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	logoutRouter.allowMethods(http.MethodPost)
	logoutRouter.preflight()
	logoutRouter.handler = &logoutHandler{}
	http.Handle("/logout", &logoutRouter)
	// http.Handle("/meetings/", cors(&meetingsHandler{}))
	meetingsRouter := router{}
	meetingsRouter.cors()
	meetingsRouter.handler = &meetingsHandler{}
	http.Handle("/meetings/", &meetingsRouter)
	// http.Handle("/meeting_messages", cors(&meetingMessagesHandler{}))
	meetingMessagesRouter := router{}
	meetingMessagesRouter.cors()
	meetingMessagesRouter.handler = &meetingMessagesHandler{}
	http.Handle("/meeting_messages/", &meetingMessagesRouter)
	// http.Handle("/meeting_users", cors(&meetingUsersHandler{}))
	meetingUsersRouter := router{}
	meetingUsersRouter.cors()
	meetingUsersRouter.handler = &meetingUsersHandler{}
	http.Handle("/meeting_users", &meetingUsersRouter)
	// http.Handle("/project_users", cors(&projectUsersHandler{}))
	projectUsersRouter := router{}
	projectUsersRouter.cors()
	projectUsersRouter.handler = &projectUsersHandler{}
	http.Handle("/project_users", &projectUsersRouter)
	// http.Handle("/projects/", cors(&projectsHandler{}))
	projectsRouter := router{}
	projectsRouter.cors()
	projectsRouter.handler = &projectsHandler{}
	http.Handle("/projects/", &projectsRouter)
	// http.Handle("/register", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodPost}, preflight(&registerHandler{}))))))
	registerRouter := router{}
	registerRouter.cors()
	registerRouter.allowCredentials(http.MethodPost)
	registerRouter.allowHeaders(map[string][]string{http.MethodPost: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	registerRouter.allowMethods(http.MethodPost)
	registerRouter.preflight()
	registerRouter.handler = &registerHandler{}
	http.Handle("/register", &registerRouter)
	// http.Handle("/repositories", cors(&repositoriesHandler{}))
	repositoriesRouter := router{}
	repositoriesRouter.cors()
	repositoriesRouter.handler = &repositoriesHandler{}
	http.Handle("/repositories", &repositoriesRouter)
	// http.Handle("/scans", cors(&scansHandler{}))
	scansRouter := router{}
	scansRouter.cors()
	scansRouter.handler = &scansHandler{}
	http.Handle("/scans", &scansRouter)
	// http.Handle("/teams/", cors(&teamsHandler{}))
	teamsRouter := router{}
	teamsRouter.cors()
	teamsRouter.handler = &teamsHandler{}
	http.Handle("/teams/", &teamsRouter)
	// http.Handle("/team_users", cors(&teamUsersHandler{}))
	teamUsersRouter := router{}
	teamUsersRouter.cors()
	teamUsersRouter.handler = &teamUsersHandler{}
	http.Handle("/team_users", &teamUsersRouter)
	// http.Handle("/team_user_invitation_requests", cors(allowMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, preflight(&teamUserInvitationRequestsHandler{}))))
	teamUserInvitationRequestsRouter := router{}
	teamUserInvitationRequestsRouter.cors()
	teamUserInvitationRequestsRouter.allowMethods(http.MethodGet, http.MethodPost, http.MethodDelete)
	teamUserInvitationRequestsRouter.preflight()
	teamUserInvitationRequestsRouter.handler = &teamUserInvitationRequestsHandler{}
	http.Handle("/team_user_invitation_requests", &teamUserInvitationRequestsRouter)
	// http.Handle("/team_user_invitation_request_projects", cors(allowMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, preflight(&teamUserInvitationRequestProjectsHandler{}))))
	teamUserInvitationRequestProjectsRouter := router{}
	teamUserInvitationRequestProjectsRouter.cors()
	teamUserInvitationRequestProjectsRouter.allowMethods(http.MethodGet, http.MethodPost, http.MethodDelete)
	teamUserInvitationRequestProjectsRouter.preflight()
	teamUserInvitationRequestProjectsRouter.handler = &teamUserInvitationRequestProjectsHandler{}
	http.Handle("/team_user_invitation_request_projects", &teamUserInvitationRequestProjectsRouter)
	// http.Handle("/tests/", cors(&testsHandler{}))
	testsRouter := router{}
	testsRouter.cors()
	testsRouter.handler = &testsHandler{}
	http.Handle("/tests/", &testsRouter)
	// http.Handle("/test_messages", cors(&testMessagesHandler{}))
	testMessagesRouter := router{}
	testMessagesRouter.cors()
	testMessagesRouter.handler = &testMessagesHandler{}
	http.Handle("/test_messages/", &testMessagesRouter)
	// http.Handle("/test_results/", cors(&testResultsHandler{}))
	testResultsRouter := router{}
	testResultsRouter.cors()
	testResultsRouter.handler = &testResultsHandler{}
	http.Handle("/test_results/", &testResultsRouter)
	// http.Handle("/test_statuses", cors(&testStatusesHandler{}))
	testStatusesRouter := router{}
	testStatusesRouter.cors()
	testStatusesRouter.handler = &testStatusesHandler{}
	http.Handle("/test_statuses", &testStatusesRouter)
	// http.Handle("/user", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodGet}, preflight(&userHandler{}))))))
	userRouter := router{}
	userRouter.cors()
	userRouter.allowCredentials(http.MethodGet)
	userRouter.allowHeaders(map[string][]string{http.MethodGet: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	userRouter.allowMethods()
	userRouter.preflight()
	userRouter.handler = &userHandler{}
	http.Handle("/user", &userRouter)
	// http.Handle("/users/", cors(&usersHandler{}))
	usersRouter := router{}
	usersRouter.cors()
	usersRouter.handler = &usersHandler{}
	http.Handle("/users/", &usersRouter)
	// http.Handle("/vulnerabilities", cors(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE}, allowMethods([]string{http.MethodGet, http.MethodPost}, preflight(&vulnerabilitiesHandler{})))))
	vulnerabilitiesRouter := router{}
	vulnerabilitiesRouter.cors()
	vulnerabilitiesRouter.handler = &vulnerabilitiesHandler{}
	http.Handle("/vulnerabilities", &vulnerabilitiesRouter)
	http.Handle("/public/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(serverHostAndPort, nil)
}
