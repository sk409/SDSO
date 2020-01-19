<template>
  <div>
    <v-form>
      <v-text-field
        v-model="username"
        :rules="usernameRules"
        label="ユーザ名"
      ></v-text-field>
      <v-text-field
        v-model="password"
        :rules="passwordRules"
        type="password"
        label="パスワード"
      ></v-text-field>
    </v-form>
    <div class="text-center">
      <v-btn color="accent" @click="submit">
        <v-icon left>{{ icon }}</v-icon>
        <span>{{ text }}</span>
      </v-btn>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true
    },
    text: {
      type: String,
      required: true
    }
  },
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
    submit() {
      this.$emit("submit", this.username, this.password);
    }
  }
};
</script>
