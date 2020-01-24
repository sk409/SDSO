<template>
  <div>
    <v-subheader>テスト結果一覧</v-subheader>
    <v-divider class="mb-3"></v-divider>
    <GitToolbar :new-revision="newRevision" @change-revision="fetchTests"></GitToolbar>
    <v-container>
      <v-row justify="center">
        <v-col cols="10">
          <v-card v-for="test in tests" :key="test.id" class="mb-8">
            <v-card-title :class="test.class" class="white--text">
              {{
              test.text
              }}
            </v-card-title>
            <v-card-text>
              <v-expansion-panels flat multiple :accordian="false">
                <v-expansion-panel
                  v-for="result in test.results"
                  :key="result.id"
                  :class="result.class"
                  class="my-3"
                >
                  <v-expansion-panel-header class="body-1">
                    {{
                    result.command
                    }}
                  </v-expansion-panel-header>
                  <v-expansion-panel-content>
                    <pre class="black white--text pa-2 console-output">{{
                      result.output
                    }}</pre>
                  </v-expansion-panel-content>
                </v-expansion-panel>
              </v-expansion-panels>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import GitToolbar from "@/components/GitToolbar.vue";
import mutations from "@/assets/js/mutations.js";
import { pathTestResults, pathTests, Url } from "@/assets/js/urls.js";
let socketTest = null;
let socketTestResult = null;
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
        username: this.user.name,
        projectname: project.name,
        revision
      };
      ajax.get(url.base, data).then(response => {
        this.tests = response.data.map(test => this.newTest(test));
      });
    },
    newTest(test) {
      const t = {};
      const classAndText = this.testClassAndText(test);
      t.class = classAndText.class;
      t.text = classAndText.text;
      test.results = test.results.map(testResult =>
        this.newTestResult(testResult)
      );
      return Object.assign(t, test);
    },
    newTestResult(testResult) {
      const t = {};
      t.class = this.resultClass(testResult);
      return Object.assign(t, testResult);
    },
    testClassAndText(test) {
      const running = test.results.some(
        result => result.status.text === "running"
      );
      if (running || test.steps != test.results.length) {
        return {
          class: "test-running",
          text: "RUNNING"
        };
      }
      const failed = test.results.some(
        result => result.status.text === "failed"
      );
      if (failed) {
        return {
          class: "test-failed",
          text: "FAILED"
        };
      }
      return {
        class: "test-success",
        text: "SUCCESS"
      };
    },
    resultClass(result) {
      switch (result.status.text) {
        case "running":
          return "test-label-running";
        case "failed":
          return "test-label-failed";
        case "success":
          return "test-label-success";
      }
    },
    setupSocket() {
      if (!WebSocket) {
        alert("WebSocketに対応していないブラウザです。");
        return;
      }
      const that = this;
      const urlTests = new Url(pathTests);
      socketTest = new WebSocket(urlTests.socket);
      socketTest.onmessage = function(e) {
        const branchname = that.$store.state.git.branchname;
        if (!branchname) {
          return;
        }
        const test = JSON.parse(e.data);
        if (branchname !== test.branchname) {
          return;
        }
        that.newRevision = true;
      };
      const urlTestResults = new Url(pathTestResults);
      socketTestResult = new WebSocket(urlTestResults.socket);
      socketTestResult.onmessage = function(e) {
        const testResult = that.newTestResult(JSON.parse(e.data));
        // TODO: gocase
        const testIndex = that.tests.findIndex(
          test => test.id === testResult.testID
        );
        const notFound = -1;
        if (testIndex === notFound) {
          return;
        }
        const test = that.tests[testIndex];
        const testResultIndex = test.results.findIndex(
          tr => tr.id === testResult.id
        );
        if (testResultIndex === notFound) {
          test.results.push(testResult);
        } else {
          that.$set(test.results, testResultIndex, testResult);
        }
        const testClassAndText = that.testClassAndText(test);
        test.text = testClassAndText.text;
        test.class = testClassAndText.class;
      };
    }
  }
};
</script>

<style>
.console-output {
  width: 600px;
  overflow-x: scroll;
}
.test-label-running {
  border-left: 5px solid rgb(130, 209, 226);
}

.test-label-failed {
  border-left: 5px solid rgb(220, 102, 97);
}

.test-label-success {
  border-left: 5px solid rgb(107, 197, 143);
}

.test-running {
  background: rgb(130, 209, 226);
}

.test-failed {
  background: rgb(220, 102, 97);
}

.test-success {
  background: rgb(107, 197, 143);
}
</style>
