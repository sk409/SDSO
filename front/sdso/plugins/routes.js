import Vue from "vue";

Vue.prototype.$routes = {
  dashboard: {
    base: "/dashboard",
    git: {
      commits: "/dashboard/git/commits",
      files(branchname, commitSHA1, path) {
        let route = "/dashboard/git/files/";
        if (!branchname || !commitSHA1) {
          return route;
        }
        route += `${branchname}/${commitSHA1}`;
        if (!path) {
          return route;
        }
        if (!path.startsWith("/")) {
          path = "/" + path;
        }
        return route + path;
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
