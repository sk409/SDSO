<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-12 col-lg-8 offset-lg-2">
        <div class="d-block d-lg-flex align-items-center mt-4">
          <div class="title text-center">プロジェクト一覧</div>
          <el-button
            type="primary"
            class="ml-auto d-none d-lg-inline"
            @click="transitionToProjectCreate()"
            >新規作成</el-button
          >
        </div>
        <el-divider class="m-0 my-3"></el-divider>
        <div>
          <div class="text-center" v-if="projects.length === 0">
            プロジェクトがありません
          </div>
          <div v-else>
            <div class="d-none d-lg-block">
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
                      <n-link
                        :to="$routes.projectCode(user.Name, project.Name)"
                        >{{ project.Name }}</n-link
                      >
                    </td>
                    <td>{{ project.CreatedAt }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <el-collapse class="d-block d-lg-none" v-model="activeNames">
              <el-collapse-item
                v-for="(project, index) in projects"
                :key="project.ID"
                :name="index"
              >
                <template slot="title">
                  <div class="ml-2">{{ project.Name }}</div>
                </template>
                <div class="border-top">
                  <div class="ml-2 mt-2">
                    <div>
                      <span>作成日:</span>
                      <span class="ml-2">{{ project.CreatedAt }}</span>
                    </div>
                    <div class="text-center mt-3">
                      <el-button
                        type="primary"
                        @click="transitionToProjectCode(project.Name)"
                        >詳細</el-button
                      >
                    </div>
                  </div>
                </div>
              </el-collapse-item>
            </el-collapse>
          </div>
        </div>
      </div>
    </div>
    <div class="fab d-lg-none" @click="transitionToProjectCreate">
      <i class="el-icon-plus"></i>
    </div>
  </div>
</template>

<script>
export default {
  layout: "Dashboard",
  data() {
    return {
      activeNames: [],
      projects: [],
      user: null
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
    transitionToProjectCreate() {
      this.$router.push(this.$routes.projectCreate);
    },
    transitionToProjectCode(projectName) {
      this.$router.push(this.$routes.projectCode(this.user.Name, projectName));
    }
  }
};
</script>

<style scoped>
.title {
  font-size: 1.25rem;
}

.fab {
  position: absolute;
  right: 28px;
  bottom: 36px;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #409eff;
  border-radius: 50%;
  color: white;
  box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.45);
  cursor: pointer;
}
</style>
