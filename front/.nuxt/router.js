import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _7b927ac6 = () => interopDefault(import('../pages/Login.vue' /* webpackChunkName: "pages/Login" */))
const _be5b15c2 = () => interopDefault(import('../pages/Register.vue' /* webpackChunkName: "pages/Register" */))
const _7828ee11 = () => interopDefault(import('../pages/dashboard/Projects.vue' /* webpackChunkName: "pages/dashboard/Projects" */))
const _3bf536a4 = () => interopDefault(import('../pages/project/Create.vue' /* webpackChunkName: "pages/project/Create" */))
const _82edfa1a = () => interopDefault(import('../pages/project/testing/_userName/_projectName/index.vue' /* webpackChunkName: "pages/project/testing/_userName/_projectName/index" */))
const _01a68e29 = () => interopDefault(import('../pages/project/vulnerabilities/_userName/_projectName/index.vue' /* webpackChunkName: "pages/project/vulnerabilities/_userName/_projectName/index" */))
const _16f4cb0f = () => interopDefault(import('../pages/project/code/_userName/_projectName/_.vue' /* webpackChunkName: "pages/project/code/_userName/_projectName/_" */))
const _fdc600b4 = () => interopDefault(import('../pages/index.vue' /* webpackChunkName: "pages/index" */))

Vue.use(Router)

export const routerOptions = {
  mode: 'history',
  base: decodeURI('/'),
  linkActiveClass: 'nuxt-link-active',
  linkExactActiveClass: 'nuxt-link-exact-active',
  scrollBehavior,

  routes: [{
    path: "/Login",
    component: _7b927ac6,
    name: "Login"
  }, {
    path: "/Register",
    component: _be5b15c2,
    name: "Register"
  }, {
    path: "/dashboard/Projects",
    component: _7828ee11,
    name: "dashboard-Projects"
  }, {
    path: "/project/Create",
    component: _3bf536a4,
    name: "project-Create"
  }, {
    path: "/project/testing/:userName?/:projectName?",
    component: _82edfa1a,
    name: "project-testing-userName-projectName"
  }, {
    path: "/project/vulnerabilities/:userName?/:projectName?",
    component: _01a68e29,
    name: "project-vulnerabilities-userName-projectName"
  }, {
    path: "/project/code/:userName?/:projectName?/*",
    component: _16f4cb0f,
    name: "project-code-userName-projectName-all"
  }, {
    path: "/",
    component: _fdc600b4,
    name: "index"
  }],

  fallback: false
}

export function createRouter () {
  return new Router(routerOptions)
}
