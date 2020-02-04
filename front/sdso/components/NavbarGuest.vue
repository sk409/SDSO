<template>
  <div>
    <v-navigation-drawer v-model="drawer" app>
      <v-list>
        <v-list-item
          v-for="navItem in navItems"
          :key="navItem.title"
          :to="navItem.route"
          router
        >
          <v-list-item-action>
            <v-icon>{{ navItem.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>
              {{ navItem.title }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar color="primary" app dark fixed>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
    </v-app-bar>
  </div>
</template>

<script>
export default {
  data() {
    return {
      drawer: false,
      navItems: [
        {
          title: "ダッシュボード",
          icon: "mdi-view-dashboard",
          route: this.$routes.dashboard.commits.base
        }
      ]
    };
  }
};
</script>
