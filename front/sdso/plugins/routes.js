import Vue from "vue";

Vue.prototype.$routes = {
  account: {
    notifications: "/account/notifications",
    settings: "/account/settings",
    teams: "/account/teams",
  },
  dashboard: {
    base: "/dashboard",
    commits: {
      base: "/dashboard/commits",
      show: "/dashboard/commits/show",
    },
    dast: "/dashboard/dast",
    files(path, file) {
      if (!path) {
        path = "";
      }
      return "/dashboard/" + path + "?file=" + !!file;
    },
    meetings: "/dashboard/meetings",
    members: "/dashboard/members",
    tests: "/dashboard/tests"
  },
  login: {
    base: "/login"
  },
  projects: {},
  register: {
    base: "/register"
  },
  teams: {
    create: "/teams/create",
    members(id) {
      return `/teams/${id}/members`;
    },
    projects(id) {
      return `/teams/${id}/projects`;
    },
    settings(id) {
      return `/teams/${id}/settings`;
    }
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
