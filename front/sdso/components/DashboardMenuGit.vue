<template>
  <v-list class="pa-2">
    <v-list-item
      v-for="item in items"
      :key="item.title"
      :to="item.route"
      router
    >
      <v-list-item-action>
        <v-icon>{{ item.icon }}</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>{{ item.title }}</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script>
export default {
  data() {
    return {
      items: [
        {
          title: "コミット",
          icon: "mdi-source-commit",
          route: this.$routes.dashboard.commits.base
        },
        {
          title: "ファイル",
          icon: "mdi-file-outline",
          route: this.$routes.dashboard.files()
        }
      ]
    };
  }
};
</script>
