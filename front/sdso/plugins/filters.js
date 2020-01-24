import Vue from "vue";
import {
  dateFormatter
} from "@/assets/js/utils.js";

Vue.filter("dateDefault", function (str) {
  console.log(str)
  return dateFormatter.default(str);
});
