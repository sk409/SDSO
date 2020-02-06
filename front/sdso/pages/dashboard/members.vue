<template>
  <MainView>
    <template v-slot:sidemenu>
      <v-list class="pa-2">
        <v-list-item
          v-for="sidemenuItem in sidemenuItems"
          :key="sidemenuItem.title"
          :to="sidemenuItem.route"
          router
        >
          <v-list-item-action>
            <v-icon>{{ sidemenuItem.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ sidemenuItem.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </template>
    <template v-slot:content>
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
            <div class="subtitle-1">{{ user.name }}</div>
          </v-card-text>
        </v-card>
      </div>
    </template>
  </MainView>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import MainView from "@/components/MainView.vue";
import mutations from "@/assets/js/mutations.js";
import { pathTeamUsers, pathUsers, Url } from "@/assets/js/urls.js";
let unsubscribe = null;
export default {
  layout: "dashboard",
  components: {
    MainView
  },
  data() {
    return {
      sidemenuItems: [
        {
          title: "メンバ",
          icon: "mdi-account-multiple-outline",
          route: this.$routes.dashboard.members
        }
      ],
      users: []
    };
  },
  created() {
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
