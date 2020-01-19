<template>
  <v-app light>
    <NavbarAuth></NavbarAuth>
    <v-content class="white black--text">
      <nuxt />
    </v-content>
  </v-app>
</template>

<script>
import NavbarAuth from "@/components/NavbarAuth.vue";
export default {
  middleware: "auth",
  components: {
    NavbarAuth
  }
};
</script>
