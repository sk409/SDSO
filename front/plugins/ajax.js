import Vue from "vue";
import Axios from "axios";

function url(path, data) {
  let url = `${process.env.APP_SERVER_ORIGIN}/${path}`;
  if (data) {
    url += "?";
    for (const key in data) {
      url += `${key}=${data[key]}&`;
    }
  }
  return url;
}

function params(data) {
  const params = new URLSearchParams();
  if (data) {
    for (const key in data) {
      params.append(key, data[key]);
    }
  }
  return params;
}

Vue.prototype.$ajax = {};

Vue.prototype.$ajax.get = function(path, data, config, completion) {
  Axios.get(url(path, data), config).then(response => {
    if (!completion) {
      return;
    }
    completion(response);
  });
};

Vue.prototype.$ajax.post = function(path, data, config, completion) {
  Axios.post(url(path), params(data), config).then(response => {
    if (!completion) {
      return;
    }
    completion(response);
  });
};
