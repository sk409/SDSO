<template>
  <div>
    <v-row justify="center" class="mt-10">
      <v-col cols="12" md="8">
        <v-card class="pa-4">
          <v-card-title>ログイン</v-card-title>
          <v-card-text>
            <AuthForm
              icon="mdi-login"
              text="ログイン"
              @submit="login"
            ></AuthForm>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <div class="text-center mt-5">
      <v-btn color="secondary" text :to="$routes.register.base"
        >アカウントをお持ちでない方</v-btn
      >
    </div>
  </div>
</template>

<script>
import AuthForm from "@/components/AuthForm.vue";
import ajax from "@/assets/js/ajax.js";
import { pathLogin, Url } from "@/assets/js/urls.js";
export default {
  components: {
    AuthForm
  },
  methods: {
    login(username, password) {
      const url = new Url(pathLogin);
      const data = {
        username,
        password
      };
      const config = {
        withCredentials: true
      };
      ajax.post(url.base, data, config).then(response => {
        if (response.data.ok) {
          this.$router.push(this.$routes.dashboard.git.commits);
        }
      });
    }
  }
};
</script>
