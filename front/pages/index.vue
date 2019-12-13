<template>
  <div class="container"></div>
</template>

<script>
import Logo from "~/components/Logo.vue";

export default {
  middleware: "auth",
  created() {
    this.$router.push(this.$routes.dashboardProjects);
  }
};
</script>

<style>
</style>
