export const state = () => ({
  team: null
});

export const mutations = {
  setTeam(state, team) {
    state.team = team;
  }
}
