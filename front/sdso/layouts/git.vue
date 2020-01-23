<template>
  <Dashboard :sidemenu-type="sidemenuType">
    <template v-slot:content>
      <div>
        <GitToolbar
          :branchname="branchname"
          :branchnames="branchnames"
          :revision="revision"
          :revisions="revisions"
          @select-branchname="selectBranchname"
          @select-revision="selectRevision"
        ></GitToolbar>
        <nuxt ref="contents" />
      </div>
    </template>
  </Dashboard>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import Dashboard from "@/components/Dashboard.vue";
import GitToolbar from "@/components/GitToolbar.vue";
import mutations from "@/assets/js/mutations.js";
import { mapMutations } from "vuex";
import { pathBranches, pathCommits, Url } from "@/assets/js/urls.js";
export default {
  middleware: "auth",
  components: {
    Dashboard,
    GitToolbar
  },
  data() {
    return {
      branchnames: [],
      commits: [],
      sidemenuType: "git",
      user: null
    };
  },
  computed: {
    branchname() {
      return this.$store.state.git.branchname;
    },
    revision() {
      return this.$store.state.git.revision;
    },
    revisions() {
      return this.commits ? this.commits.map(commit => commit.sha1) : [];
    }
  },
  created() {
    this.$nuxt.$on("setSidemenuType", this.setSidemenuType);
    this.$store.subscribe((mutation, state) => {
      switch (mutation.type) {
        case mutations.projects.setProject:
          this.setBranchname("master");
          this.setRevision("");
          this.fetchBranches();
          break;
        case mutations.git.setBranchname:
          this.setRevision("");
          this.fetchCommits();
          break;
      }
    });
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchBranches();
      this.fetchCommits();
    });
  },
  methods: {
    ...mapMutations({
      setBranchname: mutations.git.setBranchname,
      setRevision: mutations.git.setRevision
    }),
    fetchBranches() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathBranches);
      const data = {
        username: this.user.name,
        projectname: project.name
      };
      ajax.get(url.base, data).then(response => {
        this.branchnames = response.data;
      });
    },
    fetchCommits() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathCommits);
      const data = {
        username: this.user.name,
        projectname: project.name,
        branchname: this.branchname
      };
      ajax.get(url.base, data).then(response => {
        this.commits = response.data;
      });
    },
    selectBranchname(branchname) {
      this.setBranchname(branchname);
    },
    selectRevision(revision) {
      this.setRevision(revision);
    },
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
