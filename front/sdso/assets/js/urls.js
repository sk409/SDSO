export const pathAuth = "auth";
export const pathBranches = "branches";
export const pathCommits = "commits/";
export const pathFiles = "files/";
export const pathLogin = "login";
export const pathProjects = "projects";
export const pathRegister = "register";
export const pathUser = "user";

export class Url {
  constructor(path) {
    this.base = process.env.serverOrigin + "/" + path;
    switch (path) {
      case pathFiles:
        this.text = this.base + "text";
    }
  }
}
