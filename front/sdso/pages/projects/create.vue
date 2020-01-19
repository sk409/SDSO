<template>
  <div>
    <v-container fluid class="mt-5">
      <v-row justify="center">
        <v-col cols="12" md="8">
          <v-card class="pa-4">
            <v-card-title>
              プロジェクト作成
            </v-card-title>
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
              <v-btn color="accent" class="mx-auto" @click="create">
                <v-icon left>mdi-plus</v-icon>
                <span>作成</span>
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, pathProjects } from "@/assets/js/urls.js";
export default {
  layout: "auth",
  data() {
    return {
      projectname: "",
      projectnameRules: [
        v => !!v || "プロジェクト名を入力してください",
        v => (v && v.length <= 128) || "128文字以内で入力してください",
        v =>
          (v && !this.projects.some(project => project.name !== v)) ||
          "同名のプロジェクトが存在しています"
      ],
      projects: [],
      user: null
    };
  },
  created() {
    this.$fetchUser().then(response => {
      this.user = response.data;
      //   const url = new Url(urlProjects);
      //   const data = {

      //   }
      //   ajax.get(url.base, )
    });
  },
  methods: {
    create() {
      const url = new Url(pathProjects);
      const data = {
        name: this.projectname,
        userId: this.user.id
      };
      ajax.post(url.base, data).then(response => {
        console.log(response);
      });
    }
  }
};
</script>
