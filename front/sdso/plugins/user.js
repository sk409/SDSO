import ajax from "@/assets/js/ajax.js";
import Vue from "vue";
import {
  Url,
  pathUser
} from "@/assets/js/urls.js";

Vue.prototype.$fetchUser = function () {
  const url = new Url(pathUser);
  return ajax.get(url.base, {}, {
    withCredentials: true
  });
}
