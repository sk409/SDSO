import Vue from "vue";

Vue.prototype.$routes = {
  dashboard: {
    base: "/dashboard",
    commits: "/dashboard/commits",
    dast: "/dashboard/dast",
    files(path, file) {
      if (!path) {
        path = "";
      }
      return "/dashboard/" + path + "?file=" + !!file;
    },
    tests: "/dashboard/tests"
  },
  login: {
    base: "/login"
  },
  projects: {
    create: "/projects/create"
  },
  register: {
    base: "/register"
  },
  tests: {
    show(id) {
      return "/tests/" + id;
    }
  },
  vulnerabilities: {
    show(id) {
      return "/vulnerabilities/" + id;
    }
  }
};
