<template>
  <v-app light>
    <NavbarGuest></NavbarGuest>
    <v-content class="white black--text">
      <nuxt />
    </v-content>
  </v-app>
</template>

<script>
import NavbarGuest from "@/components/NavbarGuest.vue";
export default {
  components: {
    NavbarGuest
  }
};
</script>
