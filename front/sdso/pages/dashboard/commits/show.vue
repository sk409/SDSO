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
              >ファイルを確認</v-btn
            >
          </v-card-title>
          <v-card-text class="pa-3">
            <pre class="body-1 overflow-x-auto">{{ commit.message }}</pre>
          </v-card-text>
        </v-card>
        <div id="editor" class="w-100 h-100 mt-5"></div>
      </v-container>
    </template>
  </MainView>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import DashboardMenuGit from "@/components/DashboardMenuGit.vue";
import GitToolbar from "@/components/GitToolbar.vue";
import MainView from "@/components/MainView.vue";
import { pathCommits, Url } from "@/assets/js/urls.js";
let editor = null;
export default {
  layout: "dashboard",
  components: {
    DashboardMenuGit,
    GitToolbar,
    MainView
  },
  data() {
    return {
      commit: null
    };
  },
  mounted() {
    this.setupAce();
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
        editor.setValue(this.commit.diff);
      });
    },
    setupAce() {
      editor = ace.edit("editor");
      // editor.$blockScrolling = Infinity;
      editor.setTheme("ace/theme/xcode");
      editor.setFontSize(20);
      editor.setReadOnly(true);
      editor.getSession().setMode("ace/mode/diff");
    }
  }
};
</script>
