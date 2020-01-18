import Vue from "vue";

Vue.prototype.$urls = {
  auth: "auth",
  branches: "branches",
  branchProtectionRules: "branch_protection_rules",
  commits: "commits/",
  commitsShow(sha1) {
    return `commits/${sha1}/show`;
  },
  files: "files/",
  filesText: "files/text",
  login: "login",
  logout: "logout",
  projects: "projects",
  register: "register",
  repositories: "repositories",
  scans: "scans",
  tests: "tests",
  testsBranch: "tests/branch",
  testResults: "test_results",
  testStatuses: "test_statuses",
  user: "user",
  users: "users",
  vulnerabilities: "vulnerabilities"
};
