<template>
  <div>
    <div v-for="(test, index) in tests" :key="test.ID">
      <div>
        {{ test.CreatedAt }}
      </div>
      <el-collapse v-model="activeNames[index]">
        <el-collapse-item
          v-for="testResult in testResults[index]"
          :key="testResult.ID"
          :name="testResult.ID"
        >
          <template slot="title">
            <span
              style="background:green;display:inline-block;width:10px;height:100%;"
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
let user = null;
export default {
  layout: "Project",
  data() {
    return {
      project: null,
      tests: [],
      testResults: [],
      activeNames: []
    };
  },
  created() {
    this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
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
          this.tests = response.data;
          response.data.forEach(test => {
            this.activeNames.push([]);
            const data = {
              test_id: test.ID
            };
            this.$ajax.get(this.$urls.testResults, data, {}, response => {
              if (response.status !== 200) {
                return;
              }
              this.testResults.push(response.data);
            });
          });
        });
      });
    });
  }
};
</script>

<style>
.output {
  background: black;
  color: white;
  padding: 1rem;
}
</style>
