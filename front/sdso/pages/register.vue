<template>
  <div>
    <v-row justify="center" class="mt-10">
      <v-col cols="12" md="8">
        <v-card class="pa-4">
          <v-card-title>登録</v-card-title>
          <v-card-text>
            <AuthForm
              icon="mdi-account-plus"
              text="登録"
              @submit="login"
            ></AuthForm>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <div class="text-center mt-5">
      <v-btn color="secondary" text :to="$routes.login.base"
        >アカウントをお持ちの方</v-btn
      >
    </div>
  </div>
</template>

<script>
import AuthForm from "@/components/AuthForm.vue";
import ajax from "@/assets/js/ajax.js";
import { pathRegister, Url } from "@/assets/js/urls.js";
export default {
  components: {
    AuthForm
  },
  methods: {
    login(username, password) {
      const url = new Url(pathRegister);
      const data = {
        username,
        password
      };
      const config = {
        withCredentials: true
      };
      ajax.post(url.base, data, config).then(response => {
        this.$router.push(this.$routes.dashboard.commits.base);
      });
    }
  }
};
</script>
