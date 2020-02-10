package main

import (
	"net/http"

	"github.com/sk409/goconst"
)

func main() {
	rootRouter := router{}
	rootRouter.gitBasicAuth()
	rootRouter.handler = &gitHandler{}
	http.Handle("/", &rootRouter)
	authRouter := router{}
	authRouter.cors()
	authRouter.allowCredentials(http.MethodGet)
	authRouter.allowHeaders(map[string][]string{http.MethodGet: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	authRouter.allowMethods(http.MethodGet)
	authRouter.preflight()
	authRouter.handler = &authHandler{}
	http.Handle("/auth", &authRouter)
	branchesRouter := router{}
	branchesRouter.cors()
	branchesRouter.handler = &branchesHandler{}
	http.Handle("/branches", &branchesRouter)
	branchProtectionRulesRouter := router{}
	branchProtectionRulesRouter.cors()
	branchProtectionRulesRouter.handler = &branchProtectionRulesHandler{}
	http.Handle("/branch_protection_rules", &branchProtectionRulesRouter)
	commitsRouter := router{}
	commitsRouter.cors()
	commitsRouter.handler = &commitsHandler{}
	http.Handle("/commits/", &commitsRouter)
	dastVulnerabilityMessagesRouter := router{}
	dastVulnerabilityMessagesRouter.cors()
	dastVulnerabilityMessagesRouter.handler = &dastVulnerabilityMessagesHandler{}
	http.Handle("/dast_vulnerability_messages/", &dastVulnerabilityMessagesRouter)
	filesRouter := router{}
	filesRouter.cors()
	filesRouter.handler = &filesHandler{}
	http.Handle("/files/", &filesRouter)
	loginRouter := router{}
	loginRouter.cors()
	loginRouter.allowCredentials(http.MethodPost)
	loginRouter.allowHeaders(map[string][]string{http.MethodPost: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	loginRouter.allowMethods(http.MethodPost)
	loginRouter.preflight()
	loginRouter.handler = &loginHandler{}
	http.Handle("/login", &loginRouter)
	logoutRouter := router{}
	logoutRouter.cors()
	logoutRouter.allowCredentials(http.MethodPost)
	logoutRouter.allowHeaders(map[string][]string{http.MethodPost: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	logoutRouter.allowMethods(http.MethodPost)
	logoutRouter.preflight()
	logoutRouter.handler = &logoutHandler{}
	http.Handle("/logout", &logoutRouter)
	meetingsRouter := router{}
	meetingsRouter.cors()
	meetingsRouter.handler = &meetingsHandler{}
	http.Handle("/meetings/", &meetingsRouter)
	meetingMessagesRouter := router{}
	meetingMessagesRouter.cors()
	meetingMessagesRouter.handler = &meetingMessagesHandler{}
	http.Handle("/meeting_messages/", &meetingMessagesRouter)
	meetingUsersRouter := router{}
	meetingUsersRouter.cors()
	meetingUsersRouter.handler = &meetingUsersHandler{}
	http.Handle("/meeting_users", &meetingUsersRouter)
	projectUsersRouter := router{}
	projectUsersRouter.cors()
	projectUsersRouter.handler = &projectUsersHandler{}
	http.Handle("/project_users", &projectUsersRouter)
	projectsRouter := router{}
	projectsRouter.cors()
	projectsRouter.handler = &projectsHandler{}
	http.Handle("/projects/", &projectsRouter)
	registerRouter := router{}
	registerRouter.cors()
	registerRouter.allowCredentials(http.MethodPost)
	registerRouter.allowHeaders(map[string][]string{http.MethodPost: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	registerRouter.allowMethods(http.MethodPost)
	registerRouter.preflight()
	registerRouter.handler = &registerHandler{}
	http.Handle("/register", &registerRouter)
	repositoriesRouter := router{}
	repositoriesRouter.cors()
	repositoriesRouter.handler = &repositoriesHandler{}
	http.Handle("/repositories", &repositoriesRouter)
	scansRouter := router{}
	scansRouter.cors()
	scansRouter.handler = &scansHandler{}
	http.Handle("/scans", &scansRouter)
	teamsRouter := router{}
	teamsRouter.cors()
	teamsRouter.handler = &teamsHandler{}
	http.Handle("/teams/", &teamsRouter)
	teamUsersRouter := router{}
	teamUsersRouter.cors()
	teamUsersRouter.handler = &teamUsersHandler{}
	http.Handle("/team_users", &teamUsersRouter)
	teamUserInvitationRequestsRouter := router{}
	teamUserInvitationRequestsRouter.cors()
	teamUserInvitationRequestsRouter.allowMethods(http.MethodGet, http.MethodPost, http.MethodDelete)
	teamUserInvitationRequestsRouter.preflight()
	teamUserInvitationRequestsRouter.handler = &teamUserInvitationRequestsHandler{}
	http.Handle("/team_user_invitation_requests", &teamUserInvitationRequestsRouter)
	teamUserInvitationRequestProjectsRouter := router{}
	teamUserInvitationRequestProjectsRouter.cors()
	teamUserInvitationRequestProjectsRouter.allowMethods(http.MethodGet, http.MethodPost, http.MethodDelete)
	teamUserInvitationRequestProjectsRouter.preflight()
	teamUserInvitationRequestProjectsRouter.handler = &teamUserInvitationRequestProjectsHandler{}
	http.Handle("/team_user_invitation_request_projects", &teamUserInvitationRequestProjectsRouter)
	testsRouter := router{}
	testsRouter.cors()
	testsRouter.handler = &testsHandler{}
	http.Handle("/tests/", &testsRouter)
	testMessagesRouter := router{}
	testMessagesRouter.cors()
	testMessagesRouter.handler = &testMessagesHandler{}
	http.Handle("/test_messages/", &testMessagesRouter)
	testResultsRouter := router{}
	testResultsRouter.cors()
	testResultsRouter.handler = &testResultsHandler{}
	http.Handle("/test_results/", &testResultsRouter)
	testStatusesRouter := router{}
	testStatusesRouter.cors()
	testStatusesRouter.handler = &testStatusesHandler{}
	http.Handle("/test_statuses", &testStatusesRouter)
	userRouter := router{}
	userRouter.cors()
	userRouter.allowCredentials(http.MethodGet)
	userRouter.allowHeaders(map[string][]string{http.MethodGet: []string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}})
	userRouter.allowMethods()
	userRouter.preflight()
	userRouter.handler = &userHandler{}
	http.Handle("/user", &userRouter)
	usersRouter := router{}
	usersRouter.cors()
	usersRouter.handler = &usersHandler{}
	http.Handle("/users/", &usersRouter)
	vulnerabilitiesRouter := router{}
	vulnerabilitiesRouter.cors()
	vulnerabilitiesRouter.handler = &vulnerabilitiesHandler{}
	http.Handle("/vulnerabilities", &vulnerabilitiesRouter)
	http.Handle("/public/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(serverHostAndPort, nil)
}
