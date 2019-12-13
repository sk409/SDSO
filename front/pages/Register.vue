<template>
  <div>
    <AuthForm class="auth-form" type="登録" @submit="register"></AuthForm>
    <div class="text-center">
      <n-link :to="$routes.login">アカウントをお持ちの方</n-link>
    </div>
  </div>
</template>

<script>
import AuthForm from "@/components/AuthForm.vue";
export default {
  name: "login",
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
.auth-form {
  width: 60%;
  margin: 2.5rem auto;
}
</style>
