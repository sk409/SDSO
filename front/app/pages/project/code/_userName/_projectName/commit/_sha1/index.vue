<template>
  <div>
    <div v-if="commit">
      <div>{{ commit.SHA1 }}</div>
      <div>{{ commit.Message }}</div>
      <pre>
        {{ commit.Diff }}
      </pre>
    </div>
  </div>
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
        userName: this.$route.params.userName,
        projectName: this.$route.params.projectName
      };
      console.log(this.$route.params);
      this.$ajax.get(
        this.$urls.commitsShow(this.$route.params.sha1),
        data,
        {},
        response => {
          this.commit = response.data;
        }
      );
    }
  }
};
</script>
