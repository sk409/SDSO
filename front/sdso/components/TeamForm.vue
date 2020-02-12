<template>
  <v-card class="pb-3">
    <v-card-title>チーム作成</v-card-title>
    <v-card-text>
      <v-form ref="form">
        <v-text-field v-model="teamname" :rules="teamnameRules" label="名前"></v-text-field>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-btn color="accent" :loading="creating" class="mx-auto" @click="create">
        <span>作成</span>
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import FormCard from "@/components/FormCard.vue";
import roles from "@/assets/js/roles.js";
import { pathTeams, pathTeamUsers, Url } from "@/assets/js/urls.js";
let user = null;
export default {
  layout: "auth",
  components: {
    FormCard
  },
  data() {
    return {
      creating: false,
      teamname: "",
      teamnameRules: [
        v => !!v || "チーム名を入力してください",
        v => v.length <= 256 || "256文字以内で入力してください",
        v => !this.teamnames.includes(v) || "このチーム名は既に使用されています"
      ],
      teams: []
    };
  },
  computed: {
    teamnames() {
      return this.teams.map(team => team.name);
    }
  },
  created() {
    this.$fetchUser()
      .then(response => {
        user = response.data;
      })
      .then(() => {
        const url = new Url(pathTeams);
        return ajax.get(url.base);
      })
      .then(response => {
        this.teams = response.data;
      });
  },
  methods: {
    create() {
      if (!this.$refs.form.validate()) {
        return;
      }
      const url = new Url(pathTeams);
      const data = {
        name: this.teamname,
        founderUserId: user.id
      };
      this.creating = true;
      ajax.post(url.base, data).then(response => {
        const team = response.data;
        this.teams.push(team);
        const url = new Url(pathTeamUsers);
        const data = {
          teamId: team.id,
          userId: user.id,
          role: roles.team.manager
        };
        ajax.post(url.base, data).then(response => {
          this.creating = false;
          this.$emit("created", team);
        });
      });
    }
  }
};
</script>
