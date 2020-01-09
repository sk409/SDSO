<template>
  <div>{{commit}}</div>
</template>

<script>
export default {
  layout: "Project",
  data() {
    return {
      commit: null
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      const data = {
        sha1: this.$route.params.sha1
      };
      this.$ajax.get(this.$urls.commitsShow, data, {}, response => {
        console.log(response);
      });
    }
  }
};
</script>