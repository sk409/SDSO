import Vue from "vue";

Vue.prototype.$routes = {
  dashboard: {
    base: "/dashboard",
    git: {
      commits: "/dashboard/git/commits",
      files(branchname, commitSHA1, path, file) {
        const r = (route) => {
          if (!file) {
            return route;
          }
          return route + "?file=true";
        };
        let route = "/dashboard/git/files/";
        if (!branchname) {
          return r(route);
        }
        route += branchname + "/";
        if (!commitSHA1) {
          return r(route);
        }
        route += commitSHA1 + "/";
        if (!path) {
          return r(route);
        }
        if (path.startsWith("/")) {
          path = path.trim("/");
        }
        return r(route + path);
      }
    }
  },
  login: {
    base: "/login",
  },
  projects: {
    create: "/projects/create"
  },
  register: {
    base: "/register",
  }
}
