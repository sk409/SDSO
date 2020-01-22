export const state = () => ({
  project: null,
});
export const mutations = {
  setProject(state, project) {
    state.project = project;
  }
}
