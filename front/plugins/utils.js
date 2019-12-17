import Vue from "vue";

Vue.filter("formatDate", function (dateString) {
  const date = new Date(dateString);
  return `${date.getFullYear()}年${date.getMonth() +
        1}月${date.getDate()}日　${date.getHours()}時${date.getMinutes()}分${date.getSeconds()}秒`;
});
