<template>
  <div class="w-75 mx-auto">
    <div class="d-flex align-items-center mt-4">
      <div class="title">プロジェクト一覧</div>
      <el-button type="primary" class="ml-auto" @click="createProject()">新規作成</el-button>
    </div>
    <el-divider class="m-0 my-3"></el-divider>
    <div>
      <el-table :data="projects" stripe>
        <el-table-column label="プロジェクト名">
          <template slot-scope="scope">
            <el-button type="success" @click="toProjectsShow(scope.row)">{{scope.row.Name}}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="作成者">
          <div>{{user ? user.Name : ""}}</div>
        </el-table-column>
        <el-table-column label="作成日" prop="CreatedAt"></el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
export default {
  layout: "Dashboard",
  middleware: "auth",
  data() {
    return {
      user: null,
      projects: []
    };
  },
  created() {
    const that = this;
    this.$ajax.get(
      this.$urls.user,
      null,
      { withCredentials: true },
      response => {
        that.user = response.data;
        const data = {
          user_id: that.user.ID
        };
        that.$ajax.get(that.$urls.projects, data, {}, response => {
          that.projects = response.data;
        });
      }
    );
  },
  methods: {
    toProjectsShow(project) {
      this.$store.commit(this.$storeMutations.projects.set, project);
      this.$router.push(this.$routes.projectsShow);
    },
    createProject() {
      this.$router.push(this.$routes.projectsCreate);
    }
  }
};
</script>

<style scoped>
.title {
  font-size: 1.25rem;
}
</style>
