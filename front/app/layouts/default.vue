<template>
  <div>
    <AppHeader id="app-header"></AppHeader>
    <nuxt />
  </div>
</template>

<script>
import AppHeader from "@/components/AppHeader.vue";
export default {
  components: {
    AppHeader
  }
};
</script>
