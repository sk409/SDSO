const namespace = {
  git: "git/",
  notifications: "notifications/",
  projects: "projects/",
  teams: "teams/"
};
const mutations = {
  git: {
    setBranchname: namespace.git + "setBranchname",
    setRevision: namespace.git + "setRevision",
  },
  notifications: {
    setMessage: namespace.notifications + "setMessage",
  },
  projects: {
    setProject: namespace.projects + "setProject"
  },
  teams: {
    setTeam: namespace.teams + "setTeam"
  }
};
export default mutations;
