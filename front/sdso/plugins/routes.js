import Vue from "vue";

Vue.prototype.$routes = {
  account: {},
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
    members: "/dashboard/members",
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
  teams: {
    create: "/teams/create"
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
