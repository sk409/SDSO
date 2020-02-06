import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _ae7f4066 = () => interopDefault(import('../pages/login.vue' /* webpackChunkName: "pages/login" */))
const _edc4ec22 = () => interopDefault(import('../pages/register.vue' /* webpackChunkName: "pages/register" */))
const _000ed26a = () => interopDefault(import('../pages/account/notifications.vue' /* webpackChunkName: "pages/account/notifications" */))
const _53fc8b90 = () => interopDefault(import('../pages/account/teams.vue' /* webpackChunkName: "pages/account/teams" */))
const _90aefa70 = () => interopDefault(import('../pages/dashboard/commits/index.vue' /* webpackChunkName: "pages/dashboard/commits/index" */))
const _00f9fb85 = () => interopDefault(import('../pages/dashboard/dast.vue' /* webpackChunkName: "pages/dashboard/dast" */))
const _244fbcff = () => interopDefault(import('../pages/dashboard/meetings.vue' /* webpackChunkName: "pages/dashboard/meetings" */))
const _3eecd882 = () => interopDefault(import('../pages/dashboard/members.vue' /* webpackChunkName: "pages/dashboard/members" */))
const _05bcb9ca = () => interopDefault(import('../pages/dashboard/tests.vue' /* webpackChunkName: "pages/dashboard/tests" */))
const _6afea7c1 = () => interopDefault(import('../pages/teams/create.vue' /* webpackChunkName: "pages/teams/create" */))
const _ae363a12 = () => interopDefault(import('../pages/dashboard/commits/show.vue' /* webpackChunkName: "pages/dashboard/commits/show" */))
const _0ef4af9a = () => interopDefault(import('../pages/tests/_id/index.vue' /* webpackChunkName: "pages/tests/_id/index" */))
const _032c796c = () => interopDefault(import('../pages/vulnerabilities/_id/index.vue' /* webpackChunkName: "pages/vulnerabilities/_id/index" */))
const _5a2f704f = () => interopDefault(import('../pages/teams/_id/members.vue' /* webpackChunkName: "pages/teams/_id/members" */))
const _0909f8b4 = () => interopDefault(import('../pages/teams/_id/projects.vue' /* webpackChunkName: "pages/teams/_id/projects" */))
const _2050c046 = () => interopDefault(import('../pages/teams/_id/settings.vue' /* webpackChunkName: "pages/teams/_id/settings" */))
const _a446e4f0 = () => interopDefault(import('../pages/dashboard/_.vue' /* webpackChunkName: "pages/dashboard/_" */))
const _41ce8694 = () => interopDefault(import('../pages/index.vue' /* webpackChunkName: "pages/index" */))

// TODO: remove in Nuxt 3
const emptyFn = () => {}
const originalPush = Router.prototype.push
Router.prototype.push = function push (location, onComplete = emptyFn, onAbort) {
  return originalPush.call(this, location, onComplete, onAbort)
}

Vue.use(Router)

export const routerOptions = {
  mode: 'history',
  base: decodeURI('/'),
  linkActiveClass: 'nuxt-link-active',
  linkExactActiveClass: 'nuxt-link-exact-active',
  scrollBehavior,

  routes: [{
    path: "/login",
    component: _ae7f4066,
    name: "login"
  }, {
    path: "/register",
    component: _edc4ec22,
    name: "register"
  }, {
    path: "/account/notifications",
    component: _000ed26a,
    name: "account-notifications"
  }, {
    path: "/account/teams",
    component: _53fc8b90,
    name: "account-teams"
  }, {
    path: "/dashboard/commits",
    component: _90aefa70,
    name: "dashboard-commits"
  }, {
    path: "/dashboard/dast",
    component: _00f9fb85,
    name: "dashboard-dast"
  }, {
    path: "/dashboard/meetings",
    component: _244fbcff,
    name: "dashboard-meetings"
  }, {
    path: "/dashboard/members",
    component: _3eecd882,
    name: "dashboard-members"
  }, {
    path: "/dashboard/tests",
    component: _05bcb9ca,
    name: "dashboard-tests"
  }, {
    path: "/teams/create",
    component: _6afea7c1,
    name: "teams-create"
  }, {
    path: "/dashboard/commits/show",
    component: _ae363a12,
    name: "dashboard-commits-show"
  }, {
    path: "/tests/:id?",
    component: _0ef4af9a,
    name: "tests-id"
  }, {
    path: "/vulnerabilities/:id?",
    component: _032c796c,
    name: "vulnerabilities-id"
  }, {
    path: "/teams/:id?/members",
    component: _5a2f704f,
    name: "teams-id-members"
  }, {
    path: "/teams/:id?/projects",
    component: _0909f8b4,
    name: "teams-id-projects"
  }, {
    path: "/teams/:id?/settings",
    component: _2050c046,
    name: "teams-id-settings"
  }, {
    path: "/dashboard/*",
    component: _a446e4f0,
    name: "dashboard-all"
  }, {
    path: "/",
    component: _41ce8694,
    name: "index"
  }],

  fallback: false
}

export function createRouter () {
  return new Router(routerOptions)
}
