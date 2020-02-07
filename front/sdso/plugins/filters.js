import Vue from "vue";
import {
  dateFormatter,
  truncate
} from "@/assets/js/utils.js";

Vue.filter("dateAgo", function (str) {
  return dateFormatter.ago(str);
});

Vue.filter("dateDefault", function (str) {
  return dateFormatter.default(str);
});

Vue.filter("truncate", function (str, length) {
  return truncate(str, length);
});

Vue.filter("role", function (str) {
  switch (str) {
    case "manager":
      return "管理者";
    case "member":
      return "一般ユーザ";
  }
});
