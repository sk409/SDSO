<template>
  <div>
    <v-container fluid class="mt-5">
      <v-row justify="center">
        <v-col cols="12" md="8">
          <v-card class="pa-4">
            <v-card-title>プロジェクト作成</v-card-title>
            <v-card-text>
              <v-form>
                <v-text-field
                  v-model="projectname"
                  :rules="projectnameRules"
                  label="プロジェクト名"
                ></v-text-field>
              </v-form>
            </v-card-text>
            <v-card-actions>
              <v-btn
                color="accent"
                :loading="creating"
                class="mx-auto"
                @click="create"
              >
                <v-icon left>mdi-plus</v-icon>
                <span>作成</span>
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    <v-snackbar v-model="snackbar" :timeout="3000" top>
      <span>{{ notification }}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { pathProjects, Url } from "@/assets/js/urls.js";
export default {
  layout: "auth",
  data() {
    return {
      creating: false,
      notification: "",
      projectname: "",
      projectnameRules: [
        v => !!v || "プロジェクト名を入力してください",
        v => (v && v.length <= 128) || "128文字以内で入力してください",
        v =>
          (v && !this.projects.some(project => project.name === v)) ||
          "同名のプロジェクトが存在しています"
      ],
      projects: [],
      snackbar: false,
      user: null
    };
  },
  created() {
    this.$fetchUser()
      .then(async response => {
        this.user = response.data;
        const url = new Url(pathProjects);
        const data = {
          userId: this.user.id
        };
        return ajax.get(url.base, data);
      })
      .then(response => {
        this.projects = response.data;
      });
  },
  methods: {
    create() {
      const url = new Url(pathProjects);
      const data = {
        name: this.projectname,
        userId: this.user.id
      };
      this.creating = true;
      ajax.post(url.base, data).then(response => {
        this.creating = false;
        this.snackbar = true;
        if (response.status === 200) {
          this.notification = this.projectname + "を作成しました";
          this.projects.push(response.data);
        } else {
          this.notification = this.projectname + "の作成に失敗しました";
        }
      });
    }
  }
};
</script>
