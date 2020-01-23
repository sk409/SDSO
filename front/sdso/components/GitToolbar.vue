<template>
  <v-row justify="end">
    <v-col cols="3">
      <v-select
        :items="branchnames"
        no-data-text="ブランチがありません"
        placeholder="ブランチを選択してください"
        :value="branchname"
        @input="changeBranchname"
      ></v-select>
    </v-col>
    <v-col v-show="!hideRevision" cols="3" class="mr-3">
      <v-badge color="green" content="new" inline :value="newRevision">
        <v-select
          :items="revisions"
          hide-details
          no-data-text="コミットがありません"
          placeholder="コミットを選択してください"
          :value="revision"
          @input="changeRevision"
        ></v-select>
      </v-badge>
    </v-col>
  </v-row>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import { pathBranches, pathCommits, Url } from "@/assets/js/urls.js";
import { truncate } from "@/assets/js/utils.js";
import { mapMutations } from "vuex";
let unsubscribe = null;
export default {
  props: {
    hideRevision: {
      type: Boolean,
      default: false
    },
    newRevision: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      branchnames: [],
      commits: [],
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
      return this.commits.map(commit => ({
        text: truncate(commit.message, 10) + " | " + commit.sha1,
        value: commit.sha1
      }));
    }
  },
  created() {
    unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === mutations.projects.setProject) {
        this.fetchBranches();
        this.fetchCommits();
      }
    });
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchBranches();
      this.fetchCommits();
    });
  },
  destroyed() {
    unsubscribe();
  },
  watch: {
    newRevision(newValue, oldValue) {
      console.log("WATCH");
      console.log(newValue);
      console.log(oldValue);
      if (newValue === true && oldValue === false) {
        this.fetchCommits();
      }
    }
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
    fetchCommits(completion) {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      if (!this.branchname) {
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
        if (completion) {
          completion();
        }
      });
    },
    changeBranchname(branchname) {
      this.setBranchname(branchname);
      this.fetchCommits(() => {
        if (this.commits.length === 0) {
          return;
        }
        this.setRevision(this.commits[0].sha1);
        this.$emit("change-revision", this.commits[0].sha1);
      });
      this.$emit("change-branchname", branchname);
    },
    changeRevision(revision) {
      this.setRevision(revision);
      this.$emit("change-revision", revision);
    }
  }
};
</script>
