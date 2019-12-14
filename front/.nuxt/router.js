import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _2c8bcc34 = () => interopDefault(import('../pages/Login.vue' /* webpackChunkName: "pages/Login" */))
const _63c62594 = () => interopDefault(import('../pages/Register.vue' /* webpackChunkName: "pages/Register" */))
const _b56d1530 = () => interopDefault(import('../pages/dashboard/Projects.vue' /* webpackChunkName: "pages/dashboard/Projects" */))
const _1c452f85 = () => interopDefault(import('../pages/project/Create.vue' /* webpackChunkName: "pages/project/Create" */))
const _f0d4fb74 = () => interopDefault(import('../pages/project/vulnerabilities/_projectName/index.vue' /* webpackChunkName: "pages/project/vulnerabilities/_projectName/index" */))
const _3094959e = () => interopDefault(import('../pages/project/code/_projectName/_.vue' /* webpackChunkName: "pages/project/code/_projectName/_" */))
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
    path: "/project/vulnerabilities/:projectName?",
    component: _f0d4fb74,
    name: "project-vulnerabilities-projectName"
  }, {
    path: "/project/code/:projectName?/*",
    component: _3094959e,
    name: "project-code-projectName-all"
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
