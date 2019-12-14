<template>
  <div class="w-75 mx-auto mt-3">
    <h2 class="project-name">{{projectName()}}</h2>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="Code" name="Code">
        <RepositoryFiles :files="files" @file-name-click="fileNameClicked"></RepositoryFiles>
      </el-tab-pane>
      <el-tab-pane label="Test" name="Test"></el-tab-pane>
      <el-tab-pane label="Vulnerabilities" name="Vulnerabilities">
        <div v-for="scan in scans" :key="scan.ID" class="mt-3 mb-4">
          <div>{{scan.CreatedAt | formatDate}}</div>
          <el-divider class="my-2"></el-divider>
          <div v-if="scan.vulnerabilities && scan.vulnerabilities.length == 0">{{"脆弱性は検出されませんでした。"}}</div>
          <div v-else>
            <el-collapse v-model="vulnerabilityExpandedNames">
              <el-collapse-item
                v-for="vulnerability in scan.vulnerabilities"
                :key="vulnerability.ID"
                :name="vulnerability.ID"
              >
                <template slot="title">
                  <span class="ml-2">{{vulnerability.Path}}</span>
                </template>
                <div class="px-3 py-2">
                  <div>
                    <span class="vulnerability-head">種類</span>
                    <el-divider direction="vertical"></el-divider>
                    <span>{{vulnerability.Name}}</span>
                  </div>
                  <div>
                    <span class="vulnerability-head">説明</span>
                    <el-divider direction="vertical"></el-divider>
                    <span>{{vulnerability.Description}}</span>
                  </div>
                  <div>
                    <span class="vulnerability-head">メソッド</span>
                    <el-divider direction="vertical"></el-divider>
                    <span>{{vulnerability.Method}}</span>
                  </div>
                  <el-tabs
                    type="border-card"
                    v-model="requestResponseActiveTabs[vulnerability.ID]"
                    class="mt-2"
                  >
                    <el-tab-pane label="リクエスト" name="request">
                      <pre>{{decodeURI(vulnerability.Request)}}</pre>
                    </el-tab-pane>
                    <el-tab-pane label="レスポンス" name="response">
                      <pre>{{vulnerability.Response}}</pre>
                    </el-tab-pane>
                  </el-tabs>
                </div>
              </el-collapse-item>
            </el-collapse>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
// import RepositoryFiles from "@/components/RepositoryFiles.vue";

let user = null;
export default {
  // middleware: "auth",
  // components: {
  //   RepositoryFiles
  // }
  // data() {
  //   return {
  //     activeTab: "Code",
  //     requestResponseActiveTabs: {},
  //     scans: [],
  //     vulnerabilityGroups: [],
  //     vulnerabilityExpandedNames: []
  //   };
  // },
  // created() {
  //   this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
  //     user = response.data;
  //     this.fetchFiles("");
  //   });
  //   const scanData = {
  //     project_id: this.$store.state.projects.project.ID
  //   };
  //   this.$ajax.get(this.$urls.scans, scanData, {}, response => {
  //     this.scans = response.data;
  //     this.scans.forEach(scan => {
  //       const vulnerabilityData = {
  //         scan_id: scan.ID
  //       };
  //       this.$ajax.get(
  //         this.$urls.vulnerabilities,
  //         vulnerabilityData,
  //         {},
  //         response => {
  //           scan.vulnerabilities = response.data;
  //           const that = this;
  //           scan.vulnerabilities.forEach(vulnerability => {
  //             that.requestResponseActiveTabs[vulnerability.ID] = "request";
  //           });
  //         }
  //       );
  //     });
  //     this.scans = this.scans.sort((a, b) => {
  //       return a.date < b.date ? 1 : -1;
  //     });
  //   });
  // },
  // filters: {
  //   formatDate(dateString) {
  //     const date = new Date(dateString);
  //     return `${date.getFullYear()}年${date.getMonth() +
  //       1}月${date.getDate()}日　${date.getHours()}時${date.getMinutes()}分${date.getSeconds()}秒`;
  //   }
  // }
};
</script>


<style scoped>
.vulnerability-head {
  font-weight: bold;
  display: inline-block;
  width: 100px;
}
</style>