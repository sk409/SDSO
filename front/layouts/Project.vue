<template>
  <div>
    <AppHeader></AppHeader>
    <div class="w-75 mx-auto mt-3">
      <h2 class="project-name">{{ projectName }}</h2>
      <div class="border-bottom my-3">
        <n-link :to="$routes.projectCode(projectName)">Code</n-link>
        <n-link :to="$routes.projectCode(projectName)">Test</n-link>
        <n-link :to="$routes.projectCode(projectName)">Vulnerabilities</n-link>
      </div>
      <nuxt />
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
    projectName() {
      return this.$route.params.projectName
        ? this.$route.params.projectName
        : this.$route.params.pathMatch;
    }
  }
};
</script>
