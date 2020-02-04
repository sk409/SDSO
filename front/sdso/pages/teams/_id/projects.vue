<template>
  <v-container>
    <v-subheader>チーム内のプロジェクト一覧</v-subheader>
    <v-divider></v-divider>
    <v-row>
      <v-col cols="3">
        <v-btn color="primary" outlined @click="dialog = true">新規作成</v-btn>
      </v-col>
    </v-row>
    <v-card>
      <v-simple-table>
        <thead>
          <tr>
            <th>名前</th>
            <th>メンバ数</th>
            <th>作成日</th>
            <th>参加状態</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="project in projects" :key="project.id">
            <td>{{ project.name }}</td>
            <td>{{ project.users.length }}</td>
            <td>{{ project.createdAt | dateDefault }}</td>
            <td>
              <v-chip :color="status(project).color">
                {{ status(project).text }}
              </v-chip>
            </td>
          </tr>
        </tbody>
      </v-simple-table>
    </v-card>
    <v-dialog v-model="dialog" class="w-75">
      <ProjectForm :team="team" @created="createdProject"></ProjectForm>
    </v-dialog>
    <NotificationSnackbar></NotificationSnackbar>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import NotificationSnackbar from "@/components/NotificationSnackbar.vue";
import ProjectForm from "@/components/ProjectForm.vue";
import roles from "@/assets/js/roles.js";
import {
  pathProjects,
  pathTeams,
  pathTeamUsers,
  Url
} from "@/assets/js/urls.js";
export default {
  layout: "team",
  components: {
    NotificationSnackbar,
    ProjectForm
  },
  data() {
    return {
      dialog: false,
      projects: [],
      team: null,
      user: null
    };
  },
  created() {
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchProjects();
      this.fetchTeam();
    });
  },
  methods: {
    createdProject(project) {
      this.dialog = false;
      const url = new Url(pathProjects);
      const data = {
        id: project.id
      };
      ajax.get(url.base, data).then(response => {
        this.projects.push(response.data[0]);
      });
    },
    fetchProjects() {
      const teamId = this.$route.params.id;
      const url = new Url(pathProjects);
      const data = {
        teamId
      };
      ajax.get(url.base, data).then(response => {
        this.projects = response.data;
      });
    },
    fetchTeam() {
      const teamId = this.$route.params.id;
      const url = new Url(pathTeams);
      const data = {
        id: teamId
      };
      ajax.get(url.base, data).then(response => {
        this.team = response.data[0];
      });
    },
    status(project) {
      const index = project.users.findIndex(user => user.id === this.user.id);
      const notFound = -1;
      if (index === notFound) {
        return {
          color: "primary",
          text: "未参加"
        };
      } else {
        return {
          color: "accent",
          text: "参加済み"
        };
      }
    }
  }
};
</script>
