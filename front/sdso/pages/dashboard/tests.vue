<template>
  <MainView>
    <template v-slot:sidemenu>
      <v-list class="pa-2">
        <v-list-item
          v-for="sidemenuItem in sidemenuItems"
          :key="sidemenuItem.title"
          :to="sidemenuItem.route"
          router
        >
          <v-list-item-action>
            <v-icon>{{ sidemenuItem.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ sidemenuItem.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </template>
    <template v-slot:content>
      <v-subheader>テスト結果一覧</v-subheader>
      <v-divider class="mb-1"></v-divider>
      <GitToolbar
        class="mb-3"
        :new-revision.sync="newRevision"
        @change-revision="fetchTests"
      ></GitToolbar>
      <v-row v-if="tests.length !== 0" justify="center">
        <v-col cols="11">
          <v-subheader>このコミットのテスト</v-subheader>
          <v-divider class="mb-4"></v-divider>
          <v-card class="mb-4">
            <TestsTable :tests="tests.slice(0, 1)"></TestsTable>
          </v-card>
        </v-col>
      </v-row>
      <v-row v-if="2 <= tests.length" justify="center">
        <v-col cols="11">
          <v-subheader>このコミット以前のテスト</v-subheader>
          <v-divider class="mb-4"></v-divider>
          <v-card class="mb-4">
            <TestsTable :tests="tests.slice(1)"></TestsTable>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </MainView>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import GitToolbar from "@/components/GitToolbar.vue";
import MainView from "@/components/MainView.vue";
import mutations from "@/assets/js/mutations.js";
import TestsTable from "@/components/TestsTable.vue";
import { pathTestResults, pathTests, Url } from "@/assets/js/urls.js";

let socket = null;
export default {
  layout: "dashboard",
  components: {
    GitToolbar,
    MainView,
    TestsTable
  },
  data() {
    return {
      newRevision: false,
      sidemenuItems: [
        {
          title: "テスト結果",
          icon: "mdi-test-tube",
          route: this.$routes.dashboard.tests
        }
      ],
      tests: [],
      user: null
    };
  },
  created() {
    this.setupSocket();
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchTests();
    });
  },
  methods: {
    fetchTests() {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const revision = this.$store.state.git.revision;
      if (!revision) {
        return;
      }
      const url = new Url(pathTests);
      const data = {
        teamname: team.name,
        projectname: project.name,
        revision
      };
      ajax.get(url.revision, data).then(response => {
        this.tests = response.data;
      });
    },
    setupSocket() {
      if (!WebSocket) {
        alert("WebSocketに対応していないブラウザです。");
        return;
      }
      const that = this;
      const url = new Url(pathTests);
      socket = new WebSocket(url.socket);
      socket.onmessage = function(e) {
        const branchname = that.$store.state.git.branchname;
        if (!branchname) {
          return;
        }
        const test = JSON.parse(e.data);
        if (branchname !== test.branchname) {
          return;
        }
        const index = that.tests.findIndex(t => t.id === test.id);
        const notFound = -1;
        if (index === notFound) {
          that.newRevision = true;
        } else {
          that.$set(that.tests, index, test);
        }
      };
    }
  }
};
</script>
