<template>
  <div>
    <AppHeader></AppHeader>
    <div>
      <div class="header px-3 pt-3">
        <h2 class="project-name">{{ projectName }}</h2>
        <div class="tabs">
          <n-link :to="$routes.projectCode(userName, projectName)">Code</n-link>
          <n-link :to="$routes.projectTest(userName, projectName)">Test</n-link>
          <n-link :to="$routes.projectVulnerabilities(userName, projectName)">Vulnerabilities</n-link>
        </div>
      </div>
      <div class="p-3">
        <nuxt />
      </div>
    </div>
  </div>
</template>

<script>
import AppHeader from "@/components/AppHeader.vue";
export default {
  middleware: "auth",
  components: {
    AppHeader
  },
  computed: {
    userName() {
      return this.$route.params.userName;
    },
    projectName() {
      return this.$route.params.projectName
        ? this.$route.params.projectName
        : this.$route.params.pathMatch;
    }
  },
  methods: {}
};
</script>

<style scoped>
/* border-bottom my-3 */
/* .header {
  background: rgb(240, 240, 240);
}
.tabs {
} */
</style>
