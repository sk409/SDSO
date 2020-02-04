<template>
  <v-card class="pb-3">
    <v-card-title>
      プロジェクト作成
    </v-card-title>
    <v-card-text>
      <v-form ref="form">
        <v-text-field
          v-model="projectname"
          :rules="projectnameRules"
          label="プロジェクト名"
        ></v-text-field>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-btn color="accent" :loading="creating" class="mx-auto" @click="create">
        <v-icon left>mdi-plus</v-icon>
        <span>作成</span>
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import FormCard from "@/components/FormCard.vue";
import mutations from "@/assets/js/mutations.js";
import roles from "@/assets/js/roles.js";
import {
  pathProjects,
  pathProjectUsers,
  pathTeams,
  pathTeamUsers,
  Url
} from "@/assets/js/urls.js";
import { mapMutations } from "vuex";
let user = null;
export default {
  layout: "auth",
  props: {
    team: {
      validator: v => typeof v === "object" || v === null,
      required: true
    }
  },
  components: {
    FormCard
  },
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
      snackbar: false
    };
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
    });
  },
  methods: {
    ...mapMutations({
      setNotificationMessage: mutations.notifications.setMessage
    }),
    create() {
      if (!this.$refs.form.validate()) {
        return;
      }
      const url = new Url(pathProjects);
      const data = {
        name: this.projectname,
        teamId: this.team.id
      };
      this.creating = true;
      ajax
        .post(url.base, data)
        .then(response => {
          const project = response.data;
          const url = new Url(pathProjectUsers);
          const data = {
            projectId: project.id,
            role: roles.project.manager,
            userId: user.id
          };
          return ajax.post(url.base, data);
        })
        .then(response => {
          this.creating = false;
          this.setNotificationMessage(
            `プロジェクト「${this.projectname}」を作成しました`
          );
          this.$emit("created", response.data);
        });
    }
  }
};
</script>
