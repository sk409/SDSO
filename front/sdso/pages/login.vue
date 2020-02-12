<template>
  <div>
    <v-row justify="center" class="mt-10">
      <v-col cols="12" md="8">
        <v-card class="pa-4">
          <v-card-title>ログイン</v-card-title>
          <v-card-text>
            <v-form>
              <v-text-field v-model="username" :rules="usernameRules" label="ユーザ名"></v-text-field>
              <v-text-field
                v-model="password"
                autocomplete
                :rules="passwordRules"
                type="password"
                label="パスワード"
              ></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn color="accent" class="mx-auto" @click="login">ログイン</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
    <div class="text-center mt-5">
      <v-btn color="secondary" text :to="$routes.register.base">アカウントをお持ちでない方</v-btn>
    </div>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { pathLogin, Url } from "@/assets/js/urls.js";
export default {
  data() {
    return {
      username: "",
      usernameRules: [
        v => !!v || "ユーザ名を入力してください",
        v => (v && v.length <= 32) || "32文字以内で入力してください"
      ],
      password: "",
      passwordRules: [v => !!v || "パスワードを入力してください"]
    };
  },
  methods: {
    login() {
      const url = new Url(pathLogin);
      const data = {
        username: this.username,
        password: this.password
      };
      const config = {
        withCredentials: true
      };
      ajax.post(url.base, data, config).then(response => {
        if (response.data.ok) {
          this.$router.push(this.$routes.dashboard.commits.base);
        }
      });
    }
  }
};
</script>
