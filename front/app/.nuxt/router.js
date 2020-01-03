import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _6f14d94f = () => interopDefault(import('../pages/Login.vue' /* webpackChunkName: "pages/Login" */))
const _6985892d = () => interopDefault(import('../pages/Register.vue' /* webpackChunkName: "pages/Register" */))
const _3f94c89f = () => interopDefault(import('../pages/dashboard/Projects.vue' /* webpackChunkName: "pages/dashboard/Projects" */))
const _cd20affe = () => interopDefault(import('../pages/project/testing/_userName/_projectName/index.vue' /* webpackChunkName: "pages/project/testing/_userName/_projectName/index" */))
const _64537592 = () => interopDefault(import('../pages/project/vulnerabilities/_userName/_projectName/index.vue' /* webpackChunkName: "pages/project/vulnerabilities/_userName/_projectName/index" */))
const _5408cc7e = () => interopDefault(import('../pages/project/code/_userName/_projectName/_.vue' /* webpackChunkName: "pages/project/code/_userName/_projectName/_" */))
const _2dfb1658 = () => interopDefault(import('../pages/index.vue' /* webpackChunkName: "pages/index" */))

Vue.use(Router)

export const routerOptions = {
  mode: 'history',
  base: decodeURI('/'),
  linkActiveClass: 'nuxt-link-active',
  linkExactActiveClass: 'nuxt-link-exact-active',
  scrollBehavior,

  routes: [{
    path: "/Login",
    component: _6f14d94f,
    name: "Login"
  }, {
    path: "/Register",
    component: _6985892d,
    name: "Register"
  }, {
    path: "/dashboard/Projects",
    component: _3f94c89f,
    name: "dashboard-Projects"
  }, {
    path: "/project/testing/:userName?/:projectName?",
    component: _cd20affe,
    name: "project-testing-userName-projectName"
  }, {
    path: "/project/vulnerabilities/:userName?/:projectName?",
    component: _64537592,
    name: "project-vulnerabilities-userName-projectName"
  }, {
    path: "/project/code/:userName?/:projectName?/*",
    component: _5408cc7e,
    name: "project-code-userName-projectName-all"
  }, {
    path: "/",
    component: _2dfb1658,
    name: "index"
  }],

  fallback: false
}

export function createRouter () {
  return new Router(routerOptions)
}
