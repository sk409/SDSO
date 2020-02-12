<template>
  <div>
    <v-container>
      <v-subheader>参加しているチーム一覧</v-subheader>
      <v-divider></v-divider>
      <v-row>
        <v-col cols="3">
          <v-btn color="primary" outlined @click="dialog = true">新規作成</v-btn>
        </v-col>
      </v-row>
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
    </v-container>
    <v-dialog v-model="dialog">
      <TeamForm @created="createdTeam"></TeamForm>
    </v-dialog>
    <v-snackbar v-model="snackbar" :timeout="2000" top>
      <span>{{notification}}</span>
      <v-btn icon @click="snackbar=false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import TeamForm from "@/components/TeamForm.vue";
import { pathTeams, pathTeamUsers, Url } from "@/assets/js/urls.js";
let user = null;
export default {
  layout: "account",
  components: {
    TeamForm
  },
  data() {
    return {
      dialog: false,
      notification: "",
      snackbar: false,
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
    createdTeam(team) {
      this.dialog = false;
      this.snackbar = true;
      this.notification = `チーム「${team.name}」を作成しました`;
      this.fetchTeams();
    },
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
