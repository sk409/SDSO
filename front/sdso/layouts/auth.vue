<template>
  <v-app light>
    <NavbarAccount></NavbarAccount>
    <v-content class="white black--text">
      <nuxt />
    </v-content>
  </v-app>
</template>

<script>
import NavbarAccount from "@/components/NavbarAccount.vue";
export default {
  middleware: "auth",
  components: {
    NavbarAccount
  }
};
</script>
