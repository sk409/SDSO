<template>
  <div>
    <v-subheader>テスト結果一覧</v-subheader>
    <v-divider class="mb-1"></v-divider>
    <GitToolbar
      class="mb-3"
      :new-revision.sync="newRevision"
      @change-revision="fetchTests"
    ></GitToolbar>
    <v-row v-if="tests.length !== 0" justify="center">
      <v-col cols="11">
        <v-card class="mb-4">
          <v-simple-table>
            <thead>
              <tr>
                <th>ステータス</th>
                <th>ブランチ</th>
                <th>SHA1</th>
                <th>ステップ</th>
                <th>実施日</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="test in tests"
                :key="test.id"
                @click="$router.push($routes.tests.show(test.id))"
              >
                <td>
                  <v-chip :color="test.color" small text-color="white">
                    {{ test.status }}
                  </v-chip>
                </td>
                <td>{{ test.branchname }}</td>
                <td>{{ test.commitSha1.substr(0, 5) }}</td>
                <td>{{ test.results.length }}/{{ test.steps }}</td>
                <td>{{ test.createdAt | dateDefault }}</td>
              </tr>
            </tbody>
          </v-simple-table>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import GitToolbar from "@/components/GitToolbar.vue";
import mutations from "@/assets/js/mutations.js";
import { pathTestResults, pathTests, Url } from "@/assets/js/urls.js";

let socket = null;
export default {
  layout: "dashboard",
  components: {
    GitToolbar
  },
  data() {
    return {
      newRevision: false,
      tests: [],
      user: null
    };
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "tests");
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
