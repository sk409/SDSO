import Vue from "vue";
import { dateFormatter, truncate } from "@/assets/js/utils.js";

Vue.filter("dateDefault", function(str) {
  return dateFormatter.default(str);
});

Vue.filter("truncate", function(str, length) {
  return truncate(str, length);
});
