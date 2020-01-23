export const state = () => ({
  branchname: "",
  revision: "",
});

export const mutations = {
  setBranchname(state, branchname) {
    state.branchname = branchname;
  },
  setRevision(state, revision) {
    state.revision = revision;
  }
}
