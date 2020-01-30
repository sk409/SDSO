<template>
  <FormCard title="プロジェクト作成">
    <template v-slot:form>
      <v-form ref="form">
        <v-select
          v-model="teamname"
          :items="teamnames"
          no-data-text="チームがありません"
          :rules="teamnameRules"
          label="チーム"
          @input="selectTeamname"
        ></v-select>
        <v-text-field v-model="projectname" :rules="projectnameRules" label="プロジェクト名"></v-text-field>
      </v-form>
    </template>
    <template v-slot:buttons>
      <v-btn color="accent" :loading="creating" class="mx-auto" @click="create">
        <v-icon left>mdi-plus</v-icon>
        <span>作成</span>
      </v-btn>
    </template>
    <!-- <v-container fluid class="mt-5">
      <v-row justify="center">
        <v-col cols="12" md="8">
          <v-card class="pa-4">
            <v-card-title>プロジェクト作成</v-card-title>
            <v-card-text>
              <v-form>
                <v-text-field v-model="projectname" :rules="projectnameRules" label="プロジェクト名"></v-text-field>
              </v-form>
            </v-card-text>
            <v-card-actions>
              <v-btn color="accent" :loading="creating" class="mx-auto" @click="create">
                <v-icon left>mdi-plus</v-icon>
                <span>作成</span>
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    <v-snackbar v-model="snackbar" :timeout="3000" top>
      <span>{{ notification }}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>-->
  </FormCard>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import FormCard from "@/components/FormCard.vue";
import mutations from "@/assets/js/mutations.js";
import roles from "@/assets/js/roles.js";
import {
  pathProjects,
  pathProjectUsers,
  pathTeams,
  pathTeamUsers,
  Url
} from "@/assets/js/urls.js";
import { mapMutations } from "vuex";
let user = null;
export default {
  layout: "auth",
  components: {
    FormCard
  },
  data() {
    return {
      creating: false,
      notification: "",
      projectname: "",
      projectnameRules: [
        v => !!v || "プロジェクト名を入力してください",
        v => (v && v.length <= 128) || "128文字以内で入力してください",
        v =>
          (v && !this.projects.some(project => project.name === v)) ||
          "同名のプロジェクトが存在しています"
      ],
      projects: [],
      snackbar: false,
      team: null,
      teamname: "",
      teamnameRules: [v => !!v || "チームを選択してください"],
      teams: []
    };
  },
  computed: {
    teamnames() {
      return this.teams.map(team => team.name);
    }
  },
  created() {
    this.$fetchUser()
      .then(response => {
        user = response.data;
        const url = new Url(pathTeamUsers);
        const data = {
          userId: user.id
        };
        return ajax.get(url.base, data);
      })
      .then(response => {
        const teamIds = response.data.map(teamUser => teamUser.teamId);
        const url = new Url(pathTeams);
        const data = {
          ids: teamIds
        };
        return ajax.get(url.ids, data);
      })
      .then(response => {
        this.teams = response.data;
        this.teams = this.teams.sort((a, b) => {
          return a.name < b.name ? -1 : 1;
        });
      });
  },
  methods: {
    ...mapMutations({
      setNotificationMessage: mutations.notifications.setMessage
    }),
    create() {
      if (!this.$refs.form.validate()) {
        return;
      }
      const url = new Url(pathProjects);
      const data = {
        name: this.projectname,
        teamId: this.team.id
      };
      this.creating = true;
      ajax
        .post(url.base, data)
        .then(response => {
          const project = response.data;
          const url = new Url(pathProjectUsers);
          const data = {
            projectId: project.id,
            role: roles.project.manager,
            userId: user.id
          };
          return ajax.post(url.base, data);
        })
        .then(response => {
          this.setNotificationMessage(`${this.projectname}を作成しました`);
          this.$router.push(this.$routes.dashboard.commits);
        });
    },
    selectTeamname(teamname) {
      this.team = this.teams.find(team => team.name === teamname);
      if (!this.team) {
        return;
      }
      const url = new Url(pathProjects);
      const data = {
        teamId: this.team.id
      };
      ajax.get(url.base, data).then(response => {
        this.projects = response.data;
      });
    }
  }
};
</script>
