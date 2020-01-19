import ajax from "@/assets/js/ajax.js";
import {
  Url,
  pathAuth
} from "@/assets/js/urls.js";

export default function ({
  redirect
}) {
  const url = new Url(pathAuth);
  ajax.get(url.base, {}, {
    withCredentials: true
  }).then(response => {
    if (!response.data.authenticated) {
      redirect("/login");
    }
  });
}
