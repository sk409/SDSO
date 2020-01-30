import Vue from "vue";

Vue.prototype.$serverUrl = (path) => {
  if (!path.startsWith("/")) {
    path = "/" + path;
  }
  return process.env.serverOrigin + path;
}
