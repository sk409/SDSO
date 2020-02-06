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
export const pathMeetings = "meetings/";
export const pathMeetingMessages = "meeting_messages";
export const pathMeetingUsers = "meeting_users";
export const pathProjects = "projects/";
export const pathProjectUsers = "project_users";
export const pathRegister = "register";
export const pathScans = "scans";
export const pathTeams = "teams/";
export const pathTeamUsers = "team_users";
export const pathTestResults = "test_results";
export const pathTeamUserInvitationRequests = "team_user_invitation_requests";
export const pathTeamUserInvitationRequestProjects =
  "team_user_invitation_request_projects";
export const pathTests = "tests/";
export const pathUser = "user";
export const pathUsers = "users/";
export const pathVulnerabilities = "vulnerabilities";

export class Url {
  constructor(path) {
    this.base = process.env.serverOrigin + "/" + path;
    this.ids = this.base + "ids";
    this.show = id => {
      const delimiter = this.base.endsWith("/") ? "" : "/";
      return this.base + delimiter + id;
    };
    switch (path) {
      case pathFiles:
        this.text = this.base + "text";
        break;
      case pathMeetings:
        this.socket = socket(path);
        break;
      case pathTests:
        this.revision = this.base + "revision";
        this.socket = socket(path);
        break;
      case pathTestResults:
        this.socket = socket(path);
        break;
    }
  }
}
