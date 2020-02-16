<template>
  <v-app light v-if="user">
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
  },
  data() {
    return {
      user: null
    };
  },
  created() {
    this.$fetchUser().then(response => {
      this.user = response.data;
    });
  }
};
</script>
