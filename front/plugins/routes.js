import Vue from "vue";

Vue.prototype.$routes = {
  dashboardProjects: "/dashboard/projects",
  login: "/login",
  projectCreate: "/project/Create",
  projectCode(projectName, filePath) {
    let path = "/project/code/" + projectName;
    if (filePath) {
      path += "/" + filePath;
    }
    return path;
  },
  // projectCodeFileText: "/project/code/FileText",
  register: "/register"
};
