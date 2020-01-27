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
	http.Handle("/projects/", cors(&projectsHandler{}))
	http.Handle("/register", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodPost}, preflight(&registerHandler{}))))))
	http.Handle("/repositories", cors(&repositoriesHandler{}))
	http.Handle("/scans", cors(&scansHandler{}))
	http.Handle("/teams/", cors(&teamsHandler{}))
	http.Handle("/team_users", cors(&teamUsersHandler{}))
	http.Handle("/tests/", cors(&testsHandler{}))
	http.Handle("/test_results/", cors(&testResultsHandler{}))
	http.Handle("/test_statuses", cors(&testStatusesHandler{}))
	http.Handle("/user", cors(allowCredentials(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, allowMethods([]string{http.MethodGet}, preflight(&userHandler{}))))))
	http.Handle("/users", cors(&usersHandler{}))
	http.Handle("/vulnerabilities", cors(allowHeaders([]string{goconst.HTTP_HEADER_CONTENT_TYPE}, allowMethods([]string{http.MethodGet, http.MethodPost}, preflight(&vulnerabilitiesHandler{})))))
	http.ListenAndServe(serverHostAndPort, nil)
}
