<template>
  <div>
    <v-subheader>テスト結果一覧</v-subheader>
    <v-divider class="mb-5"></v-divider>
    <v-container>
      <v-row justify="center">
        <v-col cols="10">
          <v-card v-for="test in tests" :key="test.id" class="mb-8">
            <v-card-title :class="test.class" class="white--text">{{
              test.text
            }}</v-card-title>
            <v-card-text>
              <v-expansion-panels flat multiple :accordian="false">
                <v-expansion-panel
                  v-for="result in test.results"
                  :key="result.id"
                  :class="result.class"
                  class="my-3"
                >
                  <v-expansion-panel-header class="body-1">{{
                    result.command
                  }}</v-expansion-panel-header>
                  <v-expansion-panel-content>
                    <pre
                      class="black white--text"
                      style="width:600px;overflow-x:scroll;"
                      >{{ result.output }}</pre
                    >
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
import mutations from "@/assets/js/mutations.js";
import { pathTestResults, pathTests, Url } from "@/assets/js/urls.js";
export default {
  layout: "git",
  data() {
    return {
      tests: [],
      user: null
    };
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "tests");
    this.subscribe();
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
        response.data.forEach(test => {
          const classAndText = this.testClassAndText(test);
          test.class = classAndText.class;
          test.text = classAndText.text;
          test.results.forEach(result => {
            result.class = this.resultClass(result);
          });
        });
        this.tests = response.data;
      });
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
        const test = JSON.parse(e.data);
        if (test.branchName !== that.$store.state.project.branchName) {
          return;
        }
        test.results = [];
        that.tests.unshift(test);
      };
      socketTestResult = new WebSocket(
        "ws://" +
          process.env.APP_SERVER_HOST +
          ":" +
          process.env.APP_SERVER_PORT +
          "/test_results/socket"
      );
      socketTestResult.onmessage = function(e) {
        const testResult = JSON.parse(e.data);
        const testIndex = that.tests.findIndex(
          test => test.ID === testResult.TestID
        );
        const notFound = -1;
        if (testIndex === notFound) {
          return;
        }
        const test = that.tests[testIndex];
        if (test.BranchName !== that.$store.state.project.branchName) {
          return;
        }
        const testResultIndex = test.results.findIndex(
          tr => tr.ID === testResult.ID
        );
        if (testResultIndex === notFound) {
          test.results.push(testResult);
          return;
        }
        that.$set(test.results, testResultIndex, testResult);
      };
    },
    subscribe() {
      this.$store.subscribe((mutation, state) => {
        switch (mutation.type) {
          case mutations.projects.setProject:
          case mutations.git.setBranchname:
            this.tests = [];
            break;
          case mutations.git.setRevision:
            this.fetchTests();
            break;
        }
      });
    }
  }
};
</script>

<style>
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
