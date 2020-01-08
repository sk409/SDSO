import Vue from "vue";

Vue.prototype.$urls = {
  authCheck: "auth/check",
  branches: "branches",
  branchProtectionRules: "branch_protection_rules",
  commits: "commits",
  files: "file",
  filesText: "files/text",
  login: "login",
  logout: "logout",
  projects: "projects",
  register: "register",
  repositoriesInit: "repositories/init",
  scans: "scans",
  tests: "tests",
  testResults: "test_results",
  testStatuses: "test_statuses",
  user: "user",
  users: "users",
  vulnerabilities: "vulnerabilities",
};
