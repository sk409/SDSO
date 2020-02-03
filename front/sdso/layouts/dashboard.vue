<template>
  <v-app light>
    <NavbarProject></NavbarProject>
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
    <v-snackbar v-model="snackbar" :timeout="3000" top @input="clearNotification">
      <span>{{ this.$store.state.notifications.message }}</span>
      <v-btn left icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-app>
</template>

<script>
import mutations from "@/assets/js/mutations.js";
import MainView from "@/components/MainView.vue";
import NavbarProject from "@/components/NavbarProject.vue";
import { mapMutations } from "vuex";
export default {
  middleware: "auth",
  components: {
    MainView,
    NavbarProject
  },
  data() {
    return {
      allSidemenuItems: {
        git: [
          {
            title: "コミット",
            icon: "mdi-source-commit",
            route: this.$routes.dashboard.commits
          },
          {
            title: "ファイル",
            icon: "mdi-file-outline",
            route: this.$routes.dashboard.files()
          }
        ],
        members: [
          {
            title: "メンバ",
            icon: "mdi-account-multiple-outline",
            route: this.$routes.dashboard.members
          }
        ],
        security: [
          {
            title: "動的テスト",
            icon: "mdi-shield-plus-outline",
            route: this.$routes.dashboard.dast
          },
          {
            title: "パッケージ",
            icon: "mdi-package-variant",
            route: ""
          }
        ],
        tests: [
          {
            title: "テスト結果",
            icon: "mdi-test-tube",
            route: this.$routes.dashboard.tests
          }
        ]
      },
      sidemenuType: "git",
      snackbar: false,
      user: null
    };
  },
  computed: {
    sidemenuItems() {
      return this.allSidemenuItems[this.sidemenuType];
    }
  },
  created() {
    this.snackbar = this.$store.state.notifications.message !== "";
    this.$nuxt.$on("setSidemenuType", this.setSidemenuType);
    this.$fetchUser().then(response => {
      this.user = response.data;
    });
  },
  methods: {
    ...mapMutations({
      setNotificationMessage: mutations.notifications.setMessage
    }),
    clearNotification() {
      this.setNotificationMessage("");
    },
    setSidemenuType(sidemenuType) {
      this.sidemenuType = sidemenuType;
    }
  }
};
</script>

<style>
.sidemenu {
  border-right: 2px solid lightgrey;
}
</style>
