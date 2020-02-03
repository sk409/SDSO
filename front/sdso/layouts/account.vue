<template>
  <v-app light>
    <NavbarAccount></NavbarAccount>
    <v-content class="white black--text h-100">
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
          <nuxt />
        </template>
      </MainView>
    </v-content>
  </v-app>
</template>

<script>
import MainView from "@/components/MainView.vue";
import NavbarAccount from "@/components/NavbarAccount.vue";
export default {
  middleware: "auth",
  components: {
    MainView,
    NavbarAccount
  },
  data() {
    return {
      sidemenuItems: [
        {
          title: "チーム",
          icon: "mdi-account-multiple-outline",
          route: this.$routes.account.teams
        },
        {
          title: "通知",
          icon: "mdi-bell",
          route: this.$routes.account.notifications
        },
        {
          title: "設定",
          icon: "mdi-account",
          route: this.$routes.account.settings
        }
      ]
    };
  }
};
</script>
