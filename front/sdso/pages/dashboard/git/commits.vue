<template>
  <div></div>
</template>

<script>
export default {
  layout: "dashboard",
  created() {
    this.setSidemenuItems();
  },
  methods: {
    setSidemenuItems() {
      const sidemenuItems = [
        {
          title: "コミット",
          icon: "mdi-source-commit",
          route: "/"
        }
      ];
      this.$nuxt.$emit("setSidemenuItems", sidemenuItems);
    }
  }
};
</script>
