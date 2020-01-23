function socket(path) {
  let url = `ws://${process.env.serverHost}:${process.env.serverPort}/`;
  if (!path.endsWith("/")) {
    path += "/";
  }
  url += path + "socket";
  return url;
}

export const pathAuth = "auth";
export const pathBranches = "branches";
export const pathCommits = "commits/";
export const pathFiles = "files/";
export const pathLogin = "login";
export const pathProjects = "projects";
export const pathRegister = "register";
export const pathTestResults = "test_results";
export const pathTests = "tests/";
export const pathUser = "user";

export class Url {
  constructor(path) {
    this.base = process.env.serverOrigin + "/" + path;
    switch (path) {
      case pathFiles:
        this.text = this.base + "text";
      case pathTests:
        this.socket = socket(path);
    }
  }
}
