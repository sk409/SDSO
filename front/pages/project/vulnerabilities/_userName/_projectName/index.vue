<template>
  <div>
    <div v-for="(scan, index) in scans" :key="scan.ID">
      {{ scan.CreatedAt | formatDate }}
      <el-collapse v-model="vulnerabilityExpandedNames">
        <el-collapse-item
          v-for="vulnerability in vulnerabilities[index]"
          :key="vulnerability.ID"
          :name="vulnerability.ID"
        >
          <template slot="title">
            <span class="ml-2">{{ vulnerability.Path }}</span>
          </template>
          <div class="px-3 py-2">
            <div>
              <span class="vulnerability-head">種類</span>
              <el-divider direction="vertical"></el-divider>
              <span>{{ vulnerability.Name }}</span>
            </div>
            <div>
              <span class="vulnerability-head">説明</span>
              <el-divider direction="vertical"></el-divider>
              <span>{{ vulnerability.Description }}</span>
            </div>
            <div>
              <span class="vulnerability-head">メソッド</span>
              <el-divider direction="vertical"></el-divider>
              <span>{{ vulnerability.Method }}</span>
            </div>
            <el-tabs
              type="border-card"
              v-model="requestResponseActiveTabs[vulnerability.ID]"
              class="mt-2"
            >
              <el-tab-pane label="リクエスト" name="request">
                <pre>{{ decodeURI(vulnerability.Request) }}</pre>
              </el-tab-pane>
              <el-tab-pane label="レスポンス" name="response">
                <pre>{{ vulnerability.Response }}</pre>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script>
// import RepositoryFiles from "@/components/RepositoryFiles.vue";
let projectName = null;
let user = null;
export default {
  layout: "Project",
  data() {
    return {
      requestResponseActiveTabs: {},
      scans: [],
      vulnerabilities: [],
      vulnerabilityGroups: [],
      vulnerabilityExpandedNames: []
    };
  },
  created() {
    projectName = this.$route.params.projectName
      ? this.$route.params.projectName
      : this.$route.params.pathMatch;
    this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
      user = response.data;
      const data = {
        name: projectName,
        user_id: user.ID
      };
      this.$ajax.get(this.$urls.projects, data, {}, response => {
        const project = response.data[0];
        const scanData = {
          project_id: project.ID
        };
        this.$ajax.get(this.$urls.scans, scanData, {}, response => {
          //console.log(project);
          this.scans = response.data;
          this.scans.forEach((scan, index) => {
            const vulnerabilityData = {
              scan_id: scan.ID
            };
            this.$ajax.get(
              this.$urls.vulnerabilities,
              vulnerabilityData,
              {},
              response => {
                // console.log(response);
                const vulnerabilities = response.data;
                const that = this;
                vulnerabilities.forEach(vulnerability => {
                  that.requestResponseActiveTabs[vulnerability.ID] = "request";
                });
                this.vulnerabilities.push(vulnerabilities);
              }
            );
          });
          this.scans = this.scans.sort((a, b) => {
            return a.date < b.date ? 1 : -1;
          });
          console.log(this.scans);
        });
      });
    });
  }
};
</script>

<style scoped>
.vulnerability-head {
  font-weight: bold;
  display: inline-block;
  width: 100px;
}
</style>
