<template>
  <v-container>
    <v-subheader>参加しているチーム一覧</v-subheader>
    <v-divider></v-divider>
    <v-row justify="center">
      <v-col cols="10">
        <v-card>
          <v-simple-table>
            <thead>
              <tr>
                <th>チーム名</th>
                <th>プロジェクト数</th>
                <th>メンバー数</th>
                <th>作成日</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="team in teams"
                :key="team.id"
                @click="$router.push($routes.teams.members(team.id))"
              >
                <td>{{ team.name }}</td>
                <td>{{ team.projects.length }}</td>
                <td>{{ team.users.length }}</td>
                <td>{{ team.createdAt | dateDefault }}</td>
              </tr>
            </tbody>
          </v-simple-table>
        </v-card>
      </v-col>
    </v-row>
    <v-btn
      color="accent"
      fab
      fixed
      right
      bottom
      @click="$router.push($routes.teams.create)"
    >
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { pathTeams, pathTeamUsers, Url } from "@/assets/js/urls.js";
let user = null;
export default {
  layout: "account",
  data() {
    return {
      teams: [],
      teamUsers: []
    };
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
      this.fetchTeams();
    });
  },
  methods: {
    fetchTeams() {
      const url = new Url(pathTeamUsers);
      const data = {
        userId: user.id
      };
      ajax
        .get(url.base, data)
        .then(response => {
          const teamUsers = response.data;
          const teamIds = teamUsers.map(teamUser => teamUser.teamId);
          const url = new Url(pathTeams);
          const data = {
            ids: teamIds
          };
          return ajax.get(url.ids, data);
        })
        .then(response => {
          this.teams = response.data;
        });
    }
  }
};
</script>
