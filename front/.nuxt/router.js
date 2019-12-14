import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _2c8bcc34 = () => interopDefault(import('../pages/Login.vue' /* webpackChunkName: "pages/Login" */))
const _63c62594 = () => interopDefault(import('../pages/Register.vue' /* webpackChunkName: "pages/Register" */))
const _b56d1530 = () => interopDefault(import('../pages/dashboard/Projects.vue' /* webpackChunkName: "pages/dashboard/Projects" */))
const _1c452f85 = () => interopDefault(import('../pages/project/Create.vue' /* webpackChunkName: "pages/project/Create" */))
const _96636d74 = () => interopDefault(import('../pages/project/Show.vue' /* webpackChunkName: "pages/project/Show" */))
const _764c95c4 = () => interopDefault(import('../pages/project/code/Files.vue' /* webpackChunkName: "pages/project/code/Files" */))
const _7f527cf2 = () => interopDefault(import('../pages/project/code/FileText.vue' /* webpackChunkName: "pages/project/code/FileText" */))
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
    path: "/project/Show",
    component: _96636d74,
    name: "project-Show"
  }, {
    path: "/project/code/Files",
    component: _764c95c4,
    name: "project-code-Files"
  }, {
    path: "/project/code/FileText",
    component: _7f527cf2,
    name: "project-code-FileText"
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
