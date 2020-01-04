<template>
  <div class="container-fluid">
    <div class="row mt-5 mb-3">
      <div class="col-10 offset-1 col-lg-8 offset-lg-2 border">
        <AuthForm class="w-100" type="ログイン" @submit="login"></AuthForm>
      </div>
    </div>
    <div class="row mb-3">
      <!-- <a
        
        @click="socialLogin('google')"
      ></div>-->
      <!-- <a
        :href="socialLoginURL('google')"
        class="col-10 offset-1 col-lg-2 offset-lg-5 p-1 text-center social-login-button social-login-google"
      >Google でログイン</a> -->
    </div>
    <div class="col-12 text-center">
      <n-link :to="$routes.register">アカウントをお持ちでない方</n-link>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import AuthForm from "@/components/AuthForm.vue";
export default {
  name: "login",
  layout: "Auth",
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
            this.$notify.success({
              message: "ログインしました",
              duration: 3000
            });
            this.$router.push("/dashboard/projects");
            return;
          }
          this.$notify.error({
            message: "登録されていません",
            duration: 3000
          });
        }
      );
    },
    socialLoginURL(provider) {
      return (
        process.env.APP_SERVER_ORIGIN + `/social_login?provider=${provider}`
      );
    }
  }
};
</script>

<style scoped>
.social-login-button {
  color: white;
}

.social-login-google {
  background: #f44336;
}
</style>
