<template>
  <div>
    <v-subheader>メンバ一覧</v-subheader>
    <v-divider class="mb-1"></v-divider>
    <div class="d-grid columns-5">
      <v-card v-for="user in users" :key="user.id" class="text-center ma-5">
        <v-responsive>
          <v-avatar size="128">
            <v-img :src="$serverUrl(user.profileImagePath)"></v-img>
          </v-avatar>
        </v-responsive>
        <v-card-text>
          <div class="subtitle-1">{{user.name}}</div>
        </v-card-text>
      </v-card>
    </div>
    <v-btn color="accent" fab fixed right bottom @click="dialog = true">
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import { pathTeamUsers, pathUsers, Url } from "@/assets/js/urls.js";
let unsubscribe = null;
export default {
  layout: "dashboard",
  data() {
    return {
      users: []
    };
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "members");
    this.fetchUsers();
    this.subscribe();
  },
  destroyed() {
    unsubscribe();
  },
  methods: {
    fetchUsers() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathTeamUsers);
      const data = {
        teamId: project.teamId
      };
      ajax
        .get(url.base, data)
        .then(response => {
          const ids = response.data.map(teamUser => teamUser.userId);
          const url = new Url(pathUsers);
          const data = {
            ids
          };
          return ajax.get(url.ids, data);
        })
        .then(response => {
          this.users = response.data;
        });
    },
    subscribe() {
      unsubscribe = this.$store.subscribe(mutation => {
        if (mutation.type === mutations.projects.setProject) {
          this.fetchUsers();
        }
      });
    }
  }
};
</script>
