export const pathAuth = "auth";
export const pathLogin = "login";
export const pathProjects = "projects";
export const pathRegister = "register";
export const pathUser = "user";

export class Url {
  constructor(path) {
    this.base = process.env.serverOrigin + "/" + path;
  }
}
