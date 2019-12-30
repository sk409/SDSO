<template>
  <div class="w-75 mt-3 mx-auto">
    <el-form ref="form" :model="form" :rules="rules">
      <el-form-item label="プロジェクト名" prop="name">
        <el-input type="text" v-model="form.name"></el-input>
      </el-form-item>
      <el-form-item class="mt-3 text-center">
        <el-button type="primary" click="create" @click="create">作成</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
let user = null;
const validateUniqueProjectName = (rule, value, callback) => {
  if (user.projects.map(project => project.Name).includes(value)) {
    callback(new Error("同じ名前のプロジェクトが存在しています"));
    return;
  }
  callback();
};
export default {
  middleware: "auth",
  data() {
    return {
      form: {
        name: ""
      },
      rules: {
        name: [
          { required: true, message: "プロジェクト名を入力してください" },
          {
            validator: validateUniqueProjectName,
            message: "同じ名前のプロジェクトが存在しています"
          }
        ]
      }
    };
  },
  created() {
    this.$ajax.get(
      this.$urls.user,
      null,
      { withCredentials: true },
      response => {
        user = response.data;
        const data = {
          user_id: user.ID
        };
        this.$ajax.get(this.$urls.projects, data, {}, response => {
          console.log(response);
          user.projects = response.data;
        });
      }
    );
  },
  methods: {
    create() {
      this.$refs.form.validate(valid => {
        if (!valid) {
          return;
        }
        const data = {
          name: this.form.name,
          userID: user.ID
        };
        this.$ajax.post(this.$urls.projects, data, {}, response => {
          if (response.status === 200) {
            const data = {
              userName: user.Name,
              projectName: this.form.name
            };
            this.$ajax.post(this.$urls.repositoriesInit, data, {}, response => {
              if (response.status === 200) {
                this.$notify.success({
                  message: "プロジェクトを作成しました",
                  duration: 3000
                });
                this.$router.push(this.$routes.dashboardProjects);
              }
            });
          }
        });
      });
    }
  }
};
</script>
