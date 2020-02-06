import Vue from "vue";
import {
  dateFormatter,
  truncate
} from "@/assets/js/utils.js";

Vue.filter("dateAgo", function (str) {
  const date = new Date(str);
  const now = new Date();
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);
  if (seconds === 0) {
    return "たった今";
  }
  if (seconds < 60) {
    return seconds + "秒前";
  }
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) {
    return minutes + "分前";
  }
  const hours = Math.floor(minutes / 60);
  if (hours < 24) {
    return hours + "時間前";
  }
  const d = Math.floor(hours / 24);
  return d + "日前"
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
})
