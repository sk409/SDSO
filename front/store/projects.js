export const state = () => ({
  project: null
})

export const mutations = {
  set(state, project) {
    state.project = project;
  }
}
