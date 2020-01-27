<template>
  <v-app light>
    <NavbarAccount></NavbarAccount>
    <v-content class="white black--text h-100">
      <div ref="content" class="d-flex h-100">
        <div class="h-100 sidemenu">
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
        </div>
        <div :style="mainStyle" class="w-100 h-100 overflow-x-hidden overflow-y-auto">
          <nuxt />
        </div>
      </div>
    </v-content>
  </v-app>
</template>

<script>
import NavbarAccount from "@/components/NavbarAccount.vue";
export default {
  middleware: "auth",
  components: {
    NavbarAccount
  },
  data() {
    return {
      mainStyle: {},
      sidemenuItems: [
        {
          title: "チーム",
          icon: "mdi-account-multiple-outline",
          route: this.$routes.dashboard.commits
        },
        {
          title: "アカウント情報",
          icon: "mdi-file-outline",
          route: this.$routes.dashboard.files()
        }
      ]
    };
  },
  mounted() {
    let maxHeight = this.$refs.content.clientHeight;
    this.mainStyle = {
      "max-height": maxHeight + "px"
    };
  }
};
</script>

<style>
.sidemenu {
  border-right: 2px solid lightgrey;
}
</style>
