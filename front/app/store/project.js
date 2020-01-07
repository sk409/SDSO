export const state = () => ({
  branchName: "master"
});

export const mutations = {
  setBranchName(state, branchName) {
    state.branchName = branchName;
  }
};
