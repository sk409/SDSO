<template>
  <v-app light>
    <NavbarAccount></NavbarAccount>
    <v-content class="white black--text h-100">
      <MainView>
        <template v-slot:sidemenu>
          <v-list>
            <v-list-item
              v-for="sidemenuItem in sidemenuItems"
              :key="sidemenuItem.title"
              :to="sidemenuItem.route"
              router
            >
              <v-list-item-action>
                <v-icon>{{sidemenuItem.icon}}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>{{sidemenuItem.title}}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </template>
        <template v-slot:content>
          <nuxt v-if="owner === true" />
          <div v-else-if="owner === false"></div>
        </template>
      </MainView>
    </v-content>
  </v-app>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import MainView from "@/components/MainView.vue";
import NavbarAccount from "@/components/NavbarAccount.vue";
import { pathTeams, pathTeamUsers, Url } from "@/assets/js/urls.js";
export default {
  middleware: "auth",
  components: {
    MainView,
    NavbarAccount
  },
  data() {
    return {
      owner: null,
      team: null
    };
  },
  computed: {
    sidemenuItems() {
      return this.team
        ? [
            {
              title: "メンバ",
              icon: "mdi-account-multiple-outline",
              route: this.$routes.teams.members(this.team.id)
            },
            {
              title: "設定",
              icon: "mdi-settings-outline",
              route: this.$routes.teams.settings(this.team.id)
            }
          ]
        : null;
    }
  },
  created() {
    this.$fetchUser().then(response => {
      const user = response.data;
      const url = new Url(pathTeamUsers);
      const data = {
        userId: user.id
      };
      ajax
        .get(url.base, data)
        .then(response => {
          const teamUsers = response.data;
          const teamId = this.$route.params.id;
          const teamUser = teamUsers.find(
            teamUser => teamUser.teamId == teamId
          );
          if (!teamUser) {
            this.owner = false;
            return;
          }
          const url = new Url(pathTeams);
          const data = {
            id: teamUser.teamId
          };
          return ajax.get(url.base, data);
        })
        .then(response => {
          this.owner = true;
          this.team = response.data[0];
        });
    });
  }
};
</script>
