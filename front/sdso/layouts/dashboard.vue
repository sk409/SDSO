<template>
  <Dashboard :sidemenu-type="sidemenuType"></Dashboard>
</template>

<script>
import Dashboard from "@/components/Dashboard.vue";
export default {
  middleware: "auth",
  components: {
    Dashboard
  },
  data() {
    return {
      sidemenuType: "git"
    };
  },
  created() {
    this.$nuxt.$on("setSidemenuType", this.setSidemenuType);
  },
  methods: {
    setSidemenuType(sidemenuType) {
      this.sidemenuType = sidemenuType;
    }
  }
};
</script>

<style>
.sidemenu {
  border-right: 2px solid lightgrey;
}
</style>
