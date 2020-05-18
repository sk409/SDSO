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
	authRouter.allowCredentials()
	authRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	authRouter.allowMethods(http.MethodGet)
	authRouter.preflight()
	authRouter.handler = &authHandler{}
	http.Handle("/auth", &authRouter)

	branchesRouter := router{}
	branchesRouter.cors()
	branchesRouter.allowCredentials()
	branchesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	branchesRouter.preflight()
	branchesRouter.auth()
	branchesRouter.handler = &branchesHandler{}
	http.Handle("/branches", &branchesRouter)

	branchProtectionRulesRouter := router{}
	branchProtectionRulesRouter.cors()
	branchProtectionRulesRouter.allowCredentials()
	branchProtectionRulesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	branchProtectionRulesRouter.preflight()
	branchProtectionRulesRouter.auth()
	branchProtectionRulesRouter.handler = &branchProtectionRulesHandler{}
	http.Handle("/branch_protection_rules", &branchProtectionRulesRouter)

	commitsRouter := router{}
	commitsRouter.cors()
	commitsRouter.allowCredentials()
	commitsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	commitsRouter.preflight()
	commitsRouter.auth()
	commitsRouter.handler = &commitsHandler{}
	http.Handle("/commits/", &commitsRouter)

	dastVulnerabilityMessagesRouter := router{}
	dastVulnerabilityMessagesRouter.cors()
	dastVulnerabilityMessagesRouter.allowCredentials()
	dastVulnerabilityMessagesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	dastVulnerabilityMessagesRouter.preflight()
	dastVulnerabilityMessagesRouter.auth()
	dastVulnerabilityMessagesRouter.handler = &dastVulnerabilityMessagesHandler{}
	http.Handle("/dast_vulnerability_messages/", &dastVulnerabilityMessagesRouter)

	filesRouter := router{}
	filesRouter.cors()
	filesRouter.allowCredentials()
	filesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	filesRouter.preflight()
	filesRouter.auth()
	filesRouter.handler = &filesHandler{}
	http.Handle("/files/", &filesRouter)

	loginRouter := router{}
	loginRouter.cors()
	loginRouter.allowCredentials()
	loginRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	loginRouter.allowMethods(http.MethodPost)
	loginRouter.preflight()
	loginRouter.handler = &loginHandler{}
	http.Handle("/login", &loginRouter)

	logoutRouter := router{}
	logoutRouter.cors()
	logoutRouter.allowCredentials()
	logoutRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	logoutRouter.allowMethods(http.MethodPost)
	logoutRouter.preflight()
	logoutRouter.handler = &logoutHandler{}
	http.Handle("/logout", &logoutRouter)

	meetingsRouter := router{}
	meetingsRouter.cors()
	meetingsRouter.allowCredentials()
	meetingsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	meetingsRouter.preflight()
	meetingsRouter.auth()
	meetingsRouter.handler = &meetingsHandler{}
	http.Handle("/meetings/", &meetingsRouter)

	meetingMessagesRouter := router{}
	meetingMessagesRouter.cors()
	meetingMessagesRouter.allowCredentials()
	meetingMessagesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	meetingMessagesRouter.preflight()
	meetingMessagesRouter.auth()
	meetingMessagesRouter.handler = &meetingMessagesHandler{}
	http.Handle("/meeting_messages/", &meetingMessagesRouter)

	meetingMessageViewersRouter := router{}
	meetingMessageViewersRouter.cors()
	meetingMessageViewersRouter.allowCredentials()
	meetingMessageViewersRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	meetingMessageViewersRouter.preflight()
	meetingMessageViewersRouter.auth()
	meetingMessageViewersRouter.handler = &meetingMessageViewersHandler{}
	http.Handle("/meeting_message_viewers", &meetingMessageViewersRouter)

	meetingUsersRouter := router{}
	meetingUsersRouter.cors()
	meetingUsersRouter.allowCredentials()
	meetingUsersRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	meetingUsersRouter.preflight()
	meetingUsersRouter.auth()
	meetingUsersRouter.handler = &meetingUsersHandler{}
	http.Handle("/meeting_users", &meetingUsersRouter)

	projectUsersRouter := router{}
	projectUsersRouter.cors()
	projectUsersRouter.allowCredentials()
	projectUsersRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	projectUsersRouter.preflight()
	projectUsersRouter.auth()
	projectUsersRouter.handler = &projectUsersHandler{}
	http.Handle("/project_users", &projectUsersRouter)

	projectsRouter := router{}
	projectsRouter.cors()
	projectsRouter.allowCredentials()
	projectsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	projectsRouter.preflight()
	projectsRouter.auth()
	projectsRouter.handler = &projectsHandler{}
	http.Handle("/projects/", &projectsRouter)

	registerRouter := router{}
	registerRouter.cors()
	registerRouter.allowCredentials()
	registerRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	registerRouter.allowMethods(http.MethodPost)
	registerRouter.preflight()
	registerRouter.handler = &registerHandler{}
	http.Handle("/register", &registerRouter)

	repositoriesRouter := router{}
	repositoriesRouter.cors()
	repositoriesRouter.allowCredentials()
	repositoriesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	repositoriesRouter.preflight()
	repositoriesRouter.auth()
	repositoriesRouter.handler = &repositoriesHandler{}
	http.Handle("/repositories", &repositoriesRouter)

	scansRouter := router{}
	scansRouter.cors()
	scansRouter.allowCredentials()
	scansRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	scansRouter.preflight()
	scansRouter.auth()
	scansRouter.handler = &scansHandler{}
	http.Handle("/scans", &scansRouter)

	teamsRouter := router{}
	teamsRouter.cors()
	teamsRouter.allowCredentials()
	teamsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	teamsRouter.preflight()
	teamsRouter.auth()
	teamsRouter.handler = &teamsHandler{}
	http.Handle("/teams/", &teamsRouter)

	teamUsersRouter := router{}
	teamUsersRouter.cors()
	teamUsersRouter.allowCredentials()
	teamUsersRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	teamUsersRouter.preflight()
	teamUsersRouter.auth()
	teamUsersRouter.handler = &teamUsersHandler{}
	http.Handle("/team_users", &teamUsersRouter)

	teamUserInvitationRequestsRouter := router{}
	teamUserInvitationRequestsRouter.cors()
	teamUserInvitationRequestsRouter.allowCredentials()
	teamUserInvitationRequestsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	teamUserInvitationRequestsRouter.allowMethods(http.MethodGet, http.MethodPost, http.MethodDelete)
	teamUserInvitationRequestsRouter.preflight()
	teamUserInvitationRequestsRouter.auth()
	teamUserInvitationRequestsRouter.handler = &teamUserInvitationRequestsHandler{}
	http.Handle("/team_user_invitation_requests", &teamUserInvitationRequestsRouter)

	teamUserInvitationRequestProjectsRouter := router{}
	teamUserInvitationRequestProjectsRouter.cors()
	teamUserInvitationRequestProjectsRouter.allowCredentials()
	teamUserInvitationRequestProjectsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	teamUserInvitationRequestProjectsRouter.allowMethods(http.MethodGet, http.MethodPost, http.MethodDelete)
	teamUserInvitationRequestProjectsRouter.preflight()
	teamUserInvitationRequestProjectsRouter.auth()
	teamUserInvitationRequestProjectsRouter.handler = &teamUserInvitationRequestProjectsHandler{}
	http.Handle("/team_user_invitation_request_projects", &teamUserInvitationRequestProjectsRouter)

	testsRouter := router{}
	testsRouter.cors()
	testsRouter.allowCredentials()
	testsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	testsRouter.preflight()
	testsRouter.auth()
	testsRouter.handler = &testsHandler{}
	http.Handle("/tests/", &testsRouter)

	testMessagesRouter := router{}
	testMessagesRouter.cors()
	testMessagesRouter.allowCredentials()
	testMessagesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	testMessagesRouter.preflight()
	testMessagesRouter.auth()
	testMessagesRouter.handler = &testMessagesHandler{}
	http.Handle("/test_messages/", &testMessagesRouter)

	testMessageViewersRouter := router{}
	testMessageViewersRouter.cors()
	testMessageViewersRouter.allowCredentials()
	testMessageViewersRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	testMessageViewersRouter.preflight()
	testMessageViewersRouter.auth()
	testMessageViewersRouter.handler = &testMessageViewersHandler{}
	http.Handle("/test_message_viewers", &testMessageViewersRouter)

	testResultsRouter := router{}
	testResultsRouter.cors()
	testResultsRouter.allowCredentials()
	testResultsRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	testResultsRouter.preflight()
	testResultsRouter.auth()
	testResultsRouter.handler = &testResultsHandler{}
	http.Handle("/test_results/", &testResultsRouter)

	testStatusesRouter := router{}
	testStatusesRouter.cors()
	testStatusesRouter.allowCredentials()
	testStatusesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	testStatusesRouter.preflight()
	testStatusesRouter.auth()
	testStatusesRouter.handler = &testStatusesHandler{}
	http.Handle("/test_statuses", &testStatusesRouter)

	userRouter := router{}
	userRouter.cors()
	userRouter.allowCredentials()
	userRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	userRouter.allowMethods()
	userRouter.preflight()
	userRouter.handler = &userHandler{}
	http.Handle("/user", &userRouter)

	usersRouter := router{}
	usersRouter.cors()
	usersRouter.allowCredentials()
	usersRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	usersRouter.preflight()
	usersRouter.auth()
	usersRouter.handler = &usersHandler{}
	http.Handle("/users/", &usersRouter)

	vulnerabilitiesRouter := router{}
	vulnerabilitiesRouter.cors()
	vulnerabilitiesRouter.allowCredentials()
	vulnerabilitiesRouter.allowHeaders(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN)
	vulnerabilitiesRouter.preflight()
	vulnerabilitiesRouter.auth()
	vulnerabilitiesRouter.handler = &vulnerabilitiesHandler{}
	http.Handle("/vulnerabilities", &vulnerabilitiesRouter)

	apiScansRouter := router{}
	apiScansRouter.handler = &apiScansHandler{}
	http.Handle("/api/scans", &apiScansRouter)

	apiUsersRouter := router{}
	apiUsersRouter.handler = &apiUsersHandler{}
	http.Handle("/api/users", &apiUsersRouter)

	apiVulnerabilitiesRouter := router{}
	apiVulnerabilitiesRouter.handler = &apiVulnerabilitiesHandler{}
	http.Handle("/api/vulnerabilities", &apiVulnerabilitiesRouter)

	http.Handle("/public/", http.FileServer(http.Dir("./")))

	http.ListenAndServe(serverHostAndPort, nil)
}
