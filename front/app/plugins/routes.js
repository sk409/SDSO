import Vue from "vue";

Vue.prototype.$routes = {
  dashboardProjects: "/dashboard/projects",
  login: "/login",
  projectCreate: "/project/Create",
  projectCode(userName, projectName, filePath) {
    let path = "/project/code/" + userName + "/" + projectName;
    if (filePath) {
      path += "/" + filePath;
    }
    return path;
  },
  projectTest(useraName, projectName) {
    return "/project/testing/" + useraName + "/" + projectName;
  },
  projectVulnerabilities(userName, projectName) {
    return "/project/vulnerabilities/" + userName + "/" + projectName;
  },
  // projectCodeFileText: "/project/code/FileText",
  register: "/register"
};