<template>
  <div class="d-flex justify-end">
    <div class="mr-3">
      <v-select
        :items="branchnames"
        hide-details
        no-data-text="ブランチがありません"
        placeholder="ブランチを選択してください"
        :value="branchname"
        class="git-select"
        @input="changeBranchname"
      ></v-select>
    </div>
    <div class="mr-3">
      <v-badge color="green" content="new" inline :value="newRevision">
        <v-select
          :items="revisions"
          hide-details
          no-data-text="コミットがありません"
          placeholder="コミットを選択してください"
          :value="revision"
          class="git-select"
          @input="changeRevision"
        ></v-select>
      </v-badge>
    </div>
  </div>
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
      commits: []
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
        if (state.projects.project) {
          this.fetchBranches(() => {
            this.changeBranchname("master");
          });
        } else {
          this.changeBranchname("");
          this.changeRevision("");
        }
      }
    });
    this.fetchBranches();
    this.fetchCommits();
  },
  destroyed() {
    unsubscribe();
  },
  watch: {
    newRevision(newValue, oldValue) {
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
    fetchBranches(completion) {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathBranches);
      const data = {
        teamname: team.name,
        projectname: project.name
      };
      ajax.get(url.base, data).then(response => {
        this.branchnames = response.data;
        if (completion) {
          completion();
        }
      });
    },
    fetchCommits(completion) {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      if (!this.branchname) {
        return;
      }
      const url = new Url(pathCommits);
      const data = {
        teamname: team.name,
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
      this.$emit("update:newRevision", false);
      this.setRevision(revision);
      this.$emit("change-revision", revision);
    }
  }
};
</script>

<style>
.git-select {
  width: 200px;
}
</style>
