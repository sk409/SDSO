<template>
  <v-app light>
    <NavbarProject></NavbarProject>
    <v-content class="white black--text h-100">
      <nuxt />
    </v-content>
  </v-app>
</template>

<script>
import NavbarProject from "@/components/NavbarProject.vue";
export default {
  middleware: "auth",
  components: {
    NavbarProject
  }
};
</script>
