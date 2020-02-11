<template>
  <MainView>
    <template v-slot:sidemenu>
      <DashboardMenuGit></DashboardMenuGit>
    </template>
    <template v-slot:content>
      <v-container fluid class="h-100">
        <v-subheader>コミットの詳細</v-subheader>
        <v-divider></v-divider>
        <GitToolbar @change-revision="fetchCommit"></GitToolbar>
        <v-card v-if="commit" class="mt-5">
          <v-card-title class="secondary white--text d-flex">
            <span>{{ commit.sha1 }}</span>
            <v-btn
              color="white"
              outlined
              class="ml-auto"
              @click="$router.push($routes.dashboard.files())"
            >ファイルを確認</v-btn>
          </v-card-title>
          <v-card-text class="pa-3">
            <pre class="body-1 overflow-x-auto">{{ commit.message }}</pre>
          </v-card-text>
        </v-card>
        <TextEditor mode="ace/mode/diff" :value="diff"></TextEditor>
      </v-container>
    </template>
  </MainView>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import DashboardMenuGit from "@/components/DashboardMenuGit.vue";
import GitToolbar from "@/components/GitToolbar.vue";
import MainView from "@/components/MainView.vue";
import TextEditor from "@/components/TextEditor.vue";
import { pathCommits, Url } from "@/assets/js/urls.js";

export default {
  layout: "dashboard",
  components: {
    DashboardMenuGit,
    GitToolbar,
    MainView,
    TextEditor
  },
  data() {
    return {
      commit: null
    };
  },
  computed: {
    diff() {
      return this.commit ? this.commit.diff : "";
    }
  },
  mounted() {
    this.fetchCommit();
  },
  methods: {
    fetchCommit() {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const sha1 = this.$store.state.git.revision;
      if (!sha1) {
        return;
      }
      const url = new Url(pathCommits);
      const data = {
        teamname: team.name,
        projectname: project.name
      };
      ajax.get(url.show(sha1), data).then(response => {
        this.commit = response.data;
      });
    }
  }
};
</script>