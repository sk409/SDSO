import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _2c8bcc34 = () => interopDefault(import('../pages/Login.vue' /* webpackChunkName: "pages/Login" */))
const _63c62594 = () => interopDefault(import('../pages/Register.vue' /* webpackChunkName: "pages/Register" */))
const _b56d1530 = () => interopDefault(import('../pages/dashboard/Projects.vue' /* webpackChunkName: "pages/dashboard/Projects" */))
const _1c452f85 = () => interopDefault(import('../pages/project/Create.vue' /* webpackChunkName: "pages/project/Create" */))
const _03e908ec = () => interopDefault(import('../pages/project/testing/_userName/_projectName/index.vue' /* webpackChunkName: "pages/project/testing/_userName/_projectName/index" */))
const _278e7fc0 = () => interopDefault(import('../pages/project/vulnerabilities/_userName/_projectName/index.vue' /* webpackChunkName: "pages/project/vulnerabilities/_userName/_projectName/index" */))
const _b43e8350 = () => interopDefault(import('../pages/project/code/_userName/_projectName/_.vue' /* webpackChunkName: "pages/project/code/_userName/_projectName/_" */))
const _aebf5222 = () => interopDefault(import('../pages/index.vue' /* webpackChunkName: "pages/index" */))

Vue.use(Router)

export const routerOptions = {
  mode: 'history',
  base: decodeURI('/'),
  linkActiveClass: 'nuxt-link-active',
  linkExactActiveClass: 'nuxt-link-exact-active',
  scrollBehavior,

  routes: [{
    path: "/Login",
    component: _2c8bcc34,
    name: "Login"
  }, {
    path: "/Register",
    component: _63c62594,
    name: "Register"
  }, {
    path: "/dashboard/Projects",
    component: _b56d1530,
    name: "dashboard-Projects"
  }, {
    path: "/project/Create",
    component: _1c452f85,
    name: "project-Create"
  }, {
    path: "/project/testing/:userName?/:projectName?",
    component: _03e908ec,
    name: "project-testing-userName-projectName"
  }, {
    path: "/project/vulnerabilities/:userName?/:projectName?",
    component: _278e7fc0,
    name: "project-vulnerabilities-userName-projectName"
  }, {
    path: "/project/code/:userName?/:projectName?/*",
    component: _b43e8350,
    name: "project-code-userName-projectName-all"
  }, {
    path: "/",
    component: _aebf5222,
    name: "index"
  }],

  fallback: false
}

export function createRouter () {
  return new Router(routerOptions)
}
