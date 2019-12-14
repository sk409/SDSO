<template>
  <div class="h-100">
    <AppHeader id="app-header"></AppHeader>
    <div id="app-body" class="d-flex">
      <div class="sidebar">SIDEBAR</div>
      <div class="main">
        <nuxt />
      </div>
    </div>
  </div>
</template>

<script>
import AppHeader from "@/components/AppHeader.vue";
export default {
  name: "Dashboard",
  middleware: "auth",
  components: {
    AppHeader
  }
};
</script>

<style scoped>
.sidebar {
  width: 15%;
  height: 100%;
  background: rgb(41, 43, 48);
}

.main {
  width: 85%;
}
</style>
