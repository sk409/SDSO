<template>
  <Dashboard :sidemenu-type="sidemenuType">
    <template v-slot:content>
      <div class="h-100">
        <GitToolbar
          :branchname="branchname"
          :branchnames="branchnames"
          :hide-revision="hideRevision"
          :new-revision="newRevision"
          :revision="revision"
          :revisions="revisions"
          class="toolbar"
          @select-branchname="selectBranchname"
          @select-revision="selectRevision"
        ></GitToolbar>
        <nuxt ref="contents" class="main" />
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
    this.$nuxt.$on("fetchCommits", this.fetchCommits);
    this.$nuxt.$on("setHideRevision", this.setHideRevision);
    this.$nuxt.$on("setNewRevision", this.setNewRevision);
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
    // this.$fetchUser().then(response => {
    //   this.user = response.data;
    //   this.fetchBranches();
    //   this.fetchCommits();
    // });
  },
  methods: {
    ...mapMutations({
      setBranchname: mutations.git.setBranchname,
      setRevision: mutations.git.setRevision
    }),
    selectBranchname(branchname) {
      this.setBranchname(branchname);
    },
    selectRevision(revision) {
      if (this.revisions[0] === revision) {
        this.newRevision = false;
      }
      this.setRevision(revision);
    },
    setHideRevision(hideRevision) {
      this.hideRevision = hideRevision;
    },
    setNewRevision(newRevision) {
      this.newRevision = newRevision;
    },
    setSidemenuType(sidemenuType) {
      this.sidemenuType = sidemenuType;
    }
  }
};
</script>

<style>
.main {
  height: 85%;
}
.sidemenu {
  border-right: 2px solid lightgrey;
}
.toolbar {
  height: 15%;
}
</style>
