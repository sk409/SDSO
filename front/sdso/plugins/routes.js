import Vue from "vue";

Vue.prototype.$routes = {
  dashboard: {
    base: "/dashboard",
    git: {
      commits: "/dashboard/git/commits"
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
