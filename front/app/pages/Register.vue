<template>
  <div class="container-fluid">
    <div class="row mt-5 mb-3">
      <div class="col-10 offset-1 col-lg-8 offset-lg-2 border">
        <AuthForm class="w-100" type="登録" @submit="register"></AuthForm>
      </div>
    </div>
    <div class="row">
      <div class="col-12 text-center">
        <n-link :to="$routes.login">アカウントをお持ちの方</n-link>
      </div>
    </div>
  </div>
</template>

<script>
import AuthForm from "@/components/AuthForm.vue";
export default {
  name: "login",
  layout: "Auth",
  components: {
    AuthForm
  },
  methods: {
    register(data) {
      this.$ajax.post(
        this.$urls.register,
        data,
        { withCredentials: true },
        response => {
          if (response.status == 200) {
            this.$notify.success({ message: "登録しました", duration: 3000 });
            this.$router.push(this.$routes.dashboardProjects);
            return;
          }
          this.$notify.error({ message: "登録に失敗しました", duration: 3000 });
        }
      );
    }
  }
};
</script>

<style scoped>
</style>
