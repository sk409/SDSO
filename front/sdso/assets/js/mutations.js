const namespace = {
  git: "git/",
  projects: "projects/"
};
const mutations = {
  git: {
    setBranchname: namespace.git + "setBranchname",
    setRevision: namespace.git + "setRevision",
  },
  projects: {
    setProject: namespace.projects + "setProject"
  }
};
export default mutations;
