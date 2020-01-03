<template>
  <div>
    <el-form ref="form" :model="form" :rules="rules" class="p-1 p-lg-3">
      <el-form-item label="ユーザ名" prop="name">
        <el-input type="text" v-model="form.name"></el-input>
      </el-form-item>
      <el-form-item label="パスワード" prop="password">
        <el-input type="password" v-model="form.password"></el-input>
      </el-form-item>
      <el-form-item class="text-center">
        <el-button type="primary" @click="submit">{{ type }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "AuthForm",
  props: {
    type: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      form: {
        name: "",
        password: ""
      },
      rules: {
        name: [
          {
            required: true,
            message: "ユーザ名を入力してください。",
            trigger: "change"
          },
          {
            max: 32,
            message: "32文字以内で入力してください",
            trigger: "change"
          }
        ],
        password: [
          {
            required: true,
            message: "パスワードを入力してください",
            trigger: "change"
          }
        ]
      }
    };
  },
  methods: {
    submit() {
      this.$refs.form.validate(valid => {
        if (!valid) {
          return;
        }
        const data = {
          name: this.form.name,
          password: this.form.password
        };
        this.$emit("submit", data);
      });
    }
  }
};
</script>

<style scoped></style>
