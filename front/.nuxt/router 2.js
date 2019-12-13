import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _2c8bcc34 = () => interopDefault(import('../pages/Login.vue' /* webpackChunkName: "pages/Login" */))
const _63c62594 = () => interopDefault(import('../pages/Register.vue' /* webpackChunkName: "pages/Register" */))
const _b56d1530 = () => interopDefault(import('../pages/dashboard/Projects.vue' /* webpackChunkName: "pages/dashboard/Projects" */))
const _529decae = () => interopDefault(import('../pages/projects/Create.vue' /* webpackChunkName: "pages/projects/Create" */))
const _255386a2 = () => interopDefault(import('../pages/projects/Show.vue' /* webpackChunkName: "pages/projects/Show" */))
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
    path: "/projects/Create",
    component: _529decae,
    name: "projects-Create"
  }, {
    path: "/projects/Show",
    component: _255386a2,
    name: "projects-Show"
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
