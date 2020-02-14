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
export const pathDastVulnerabilityMessages = "dast_vulnerability_messages/";
export const pathFiles = "files/";
export const pathLogin = "login";
export const pathMeetings = "meetings/";
export const pathMeetingMessages = "meeting_messages/";
export const pathMeetingMessageViewers = "meeting_message_viewers";
export const pathMeetingUsers = "meeting_users";
export const pathProjects = "projects/";
export const pathProjectUsers = "project_users";
export const pathRegister = "register";
export const pathScans = "scans";
export const pathTeams = "teams/";
export const pathTeamUsers = "team_users";
export const pathTeamUserInvitationRequests = "team_user_invitation_requests";
export const pathTeamUserInvitationRequestProjects =
  "team_user_invitation_request_projects";
export const pathTests = "tests/";
export const pathTestMessages = "test_messages/";
export const pathTestResults = "test_results";
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
      case pathDastVulnerabilityMessages:
        this.count = this.base + "count";
        this.range = this.base + "range";
        this.socket = userId => socket(path) + `?userId=${userId}`;
        break;
      case pathFiles:
        this.text = this.base + "text";
        break;
      case pathMeetingMessages:
        this.count = this.base + "count";
        this.countNew = this.base + "count/new";
        this.getIds = this.base + "get/ids";
        this.ids = this.base + "ids";
        this.new = this.base + "new";
        this.range = this.base + "range";
        this.socket = userId => socket(path) + `?userId=${userId}`;
        break;
      case pathTestMessages:
        this.count = this.base + "count";
        this.range = this.base + "range";
        this.socket = userId => socket(path) + `?userId=${userId}`;
        break;
      case pathTests:
        this.revision = this.base + "revision";
        this.socket = userId => socket(path) + `?userId=${userId}`;
        break;
    }
  }
}
