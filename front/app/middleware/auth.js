import Axios from "axios"

export default async function ({
  redirect
}) {
  const response = await Axios.get(
    `${process.env.APP_SERVER_ORIGIN}/auth`, {
      withCredentials: true
    }
  );
  if (!response.data.authenticated) {
    redirect("/login")
  }
}
