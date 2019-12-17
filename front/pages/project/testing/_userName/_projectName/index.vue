<template>
  <div>
    <div v-for="(test, index) in tests" :key="test.ID" class="test">
      <div class="d-flex align-items-center mb-2">
        <div :style="testResultStyle(test)" class="test-status">
          {{ testResultText(test) }}
        </div>
        <div class="ml-auto">{{ test.CreatedAt | formatDate }}</div>
      </div>
      <el-collapse v-model="activeNames[index]">
        <el-collapse-item
          v-for="testResult in test.results"
          :key="testResult.ID"
          :name="testResult.ID"
        >
          <template slot="title">
            <span
              :style="markerStyle(testResult.TestStatusID)"
              class="marker"
            ></span>
            <span class="ml-2">{{ testResult.Command }}</span>
          </template>
          <div class="p-2">
            <pre class="output">{{ testResult.Output }}</pre>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script>
let socketTest = null;
let socketTestResult = null;
let user = null;

function testStatusColor(any) {
  const colors = [
    "rgb(130, 209, 226)",
    "rgb(107, 197, 143)",
    "rgb(220, 102, 97)"
  ];
  //console.log(typeof any);
  if (typeof any === "number") {
    return colors[any - 1];
  }
  if (any === "SUCCESS") {
    return colors[1];
  } else if (any === "FAILED") {
    return colors[2];
  } else {
    return colors[0];
  }
}

function testStatus(test) {
  if (!test.results) {
    return;
  }
  // TODO: StatusID
  const failedTestResultIndex = test.results.findIndex(
    result => result.TestStatusID === 3
  );
  const notFound = -1;
  if (failedTestResultIndex !== notFound) {
    return "FAILED";
  }
  // TODO: StatusID
  const runningTestResultIndex = test.results.findIndex(
    result => result.TestStatusID === 1
  );
  console.log("==============");
  console.log(test.results.length < test.Steps);
  console.log(runningTestResultIndex !== notFound);
  if (test.results.length < test.Steps || runningTestResultIndex !== notFound) {
    return "RUNNING";
  }
  return "SUCCESS";
}

export default {
  layout: "Project",
  data() {
    return {
      activeNames: [],
      project: null,
      tests: []
    };
  },
  created() {
    this.setupSocket();
    this.fetchData();
  },
  methods: {
    testResultStyle(test) {
      if (!test.results) {
        return;
      }
      return {
        background: testStatusColor(testStatus(test))
      };
    },
    testResultText(test) {
      return testStatus(test);
    },
    markerStyle(statusId) {
      return {
        background: testStatusColor(statusId)
      };
    },
    setupSocket() {
      if (!WebSocket) {
        alert("WebSocketに対応していないブラウザです。");
        return;
      }
      socketTest = new WebSocket(
        "ws://" +
          process.env.APP_SERVER_HOST +
          ":" +
          process.env.APP_SERVER_PORT +
          "/test_socket"
      );
      const that = this;
      socketTest.onmessage = function(e) {
        const test = JSON.parse(e.data);
        test.results = [];
        that.tests.unshift(test);
      };
      socketTestResult = new WebSocket(
        "ws://" +
          process.env.APP_SERVER_HOST +
          ":" +
          process.env.APP_SERVER_PORT +
          "/test_result_socket"
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
        const testResultIndex = that.tests[testIndex].results.findIndex(
          tr => tr.ID === testResult.ID
        );
        if (testResultIndex === notFound) {
          that.tests[testIndex].results.push(testResult);
          return;
        }
        that.$set(that.tests[testIndex].results, testResultIndex, testResult);
      };
    },
    fetchData() {
      this.$ajax.get(
        this.$urls.user,
        {},
        { withCredentials: true },
        response => {
          if (response.status !== 200) {
            return;
          }
          user = response.data;
          const projectName = this.$route.params.projectName;
          const data = {
            name: projectName,
            user_id: user.ID
          };
          this.$ajax.get(this.$urls.projects, data, {}, response => {
            if (response.status !== 200) {
              return;
            }
            this.project = response.data[0];
            const data = {
              project_id: this.project.ID
            };
            this.$ajax.get(this.$urls.tests, data, {}, response => {
              if (response.status !== 200) {
                return;
              }
              this.tests = response.data.sort((a, b) => {
                return a.CreatedAt < b.CreatedAt ? 1 : -1;
              });
              response.data.forEach((test, index) => {
                this.activeNames.push([]);
                const data = {
                  test_id: test.ID
                };
                this.$ajax.get(this.$urls.testResults, data, {}, response => {
                  if (response.status !== 200) {
                    return;
                  }
                  this.$set(this.tests[index], "results", response.data);
                });
              });
            });
          });
        }
      );
    }
  }
};
</script>

<style>
.test {
  padding: 0.8rem;
  margin-bottom: 1.5rem;
  background: rgb(250, 250, 250);
  border: 1px solid lightgray;
}
.test-status {
  color: white;
  padding: 0.25rem;
}
.marker {
  display: inline-block;
  width: 10px;
  height: 100%;
}
.output {
  background: black;
  color: white;
  padding: 1rem;
}
</style>
