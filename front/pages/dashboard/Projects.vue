<template>
  <div class="w-75 mx-auto">
    <div class="d-flex align-items-center mt-4">
      <div class="title">プロジェクト一覧</div>
      <el-button type="primary" class="ml-auto" @click="createProject()">新規作成</el-button>
    </div>
    <el-divider class="m-0 my-3"></el-divider>
    <div>
      <table class="table table-border">
        <thead>
          <tr>
            <th>プロジェクト名</th>
            <th>作成日</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="project in projects" :key="project.ID">
            <td>
              <n-link :to="$routes.projectCode(user.Name, project.Name)">{{project.Name}}</n-link>
            </td>
            <td>{{project.CreatedAt}}</td>
          </tr>
        </tbody>
      </table>
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
    createProject() {
      this.$router.push(this.$routes.projectCreate);
    }
  }
};
</script>

<style scoped>
.title {
  font-size: 1.25rem;
}
</style>
