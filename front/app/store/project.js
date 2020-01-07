export const state = () => ({
  branch: "master"
})

export const mutations = {
  setBranch(state, branch) {
    state.branch = branch;
  }
}
