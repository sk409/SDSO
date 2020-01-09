<template>
  <div class="container-fluid">
    <div v-for="commit in commits" :key="commit.SHA1" class="mb-3 row">
      <div class="col-12 col-md-6 offset-md-3 d-flex align-items-center border p-2">
        <div>{{ commit.Message }}</div>
        <n-link
          :to="$routes.projectCodeCommit(pathParamUserName, pathParamProjectName, commit.SHA1)"
          class="ml-auto"
        >{{ commit.SHA1.substring(0, 5) }}</n-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CommitView",
  data() {
    return {
      commits: []
    };
  },
  computed: {
    pathParamUserName() {
      return this.$route.params.userName;
    },
    pathParamProjectName() {
      return this.$route.params.projectName
        ? this.$route.params.projectName
        : this.$route.params.pathMatch;
    }
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      const data = {
        userName: this.pathParamUserName,
        projectName: this.pathParamProjectName,
        branchName: this.$store.state.project.branchName
      };
      this.$ajax.get(this.$urls.commits, data, {}, response => {
        if (response.status !== 200) {
          return;
        }
        this.commits = response.data;
      });
    }
  }
};
</script>
