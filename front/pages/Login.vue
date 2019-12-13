<template>
  <div>
    <AuthForm class="auth-form" type="ログイン" @submit="login"></AuthForm>
    <div class="text-center">
      <n-link :to="$routes.register">アカウントをお持ちでない方</n-link>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import AuthForm from "@/components/AuthForm.vue";
export default {
  name: "login",
  components: {
    AuthForm
  },
  methods: {
    login(data) {
      this.$ajax.post(
        this.$urls.login,
        data,
        { withCredentials: true },
        response => {
          if (response.data.exist) {
            this.$router.push("/dashboard/projects");
            return;
          }
          this.$notify.error({
            message: "登録されていません",
            duration: 3000
          });
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
