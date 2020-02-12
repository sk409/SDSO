<template>
  <div>
    <v-row justify="center" class="mt-10">
      <v-col cols="12" md="8">
        <v-card class="pa-4">
          <v-card-title>登録</v-card-title>
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
              <v-text-field v-model="handlename" :rules="handlenameRules" label="ハンドルネーム"></v-text-field>
              <v-text-field v-model="email" :rules="emailRules" type="email" label="メールアドレス"></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn color="accent" class="mx-auto" @click="register">登録</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
    <div class="text-center mt-5">
      <v-btn color="secondary" text :to="$routes.login.base">アカウントをお持ちの方</v-btn>
    </div>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { pathRegister, Url } from "@/assets/js/urls.js";
export default {
  data() {
    return {
      email: "",
      emailRules: [
        v => !!v || "メールアドレスを入力してください",
        v => (v && v.length <= 254) || "254文字以内で入力してください",
        v => {
          if (v.length > 0) {
            const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return pattern.test(v) || "メールアドレスが無効です";
          }
          return true;
        }
      ],
      handlename: "",
      handlenameRules: [
        v => !!v || "ハンドルネームを入力してください",
        v => (v && v.length <= 32) || "32文字以内で入力してください"
      ],
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
    register() {
      const url = new Url(pathRegister);
      const data = {
        username: this.username,
        password: this.password,
        handlename: this.handlename,
        email: this.email
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
