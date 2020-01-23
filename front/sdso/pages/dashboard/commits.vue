<template>
  <div>
    <v-subheader>コミット一覧</v-subheader>
    <v-divider class="mb-5"></v-divider>
    <v-data-table
      :headers="tableHeaders"
      :items="commits"
      :items-per-page="20"
      no-data-text="まだコミットされていません"
      @click:row="clickCommit"
    >
    </v-data-table>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import { pathCommits, Url } from "@/assets/js/urls.js";
import { dateFormatter, truncate } from "@/assets/js/utils.js";
import { mapMutations } from "vuex";
export default {
  layout: "dashboard",
  data() {
    return {
      branchname: "master",
      commits: [],
      tableHeaders: [
        { text: "SHA1", value: "sha1" },
        { text: "コミットメッセージ", value: "message" },
        { text: "日付", value: "date" }
      ],
      user: null
    };
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "git");
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchCommits();
    });
    this.$store.subscribe((mutation, state) => {
      if (mutation.type !== mutations.projects.setProject) {
        return;
      }
      this.fetchCommits();
    });
  },
  methods: {
    ...mapMutations({
      setRevision: mutations.git.setRevision
    }),
    clickCommit(commit) {
      this.setRevision(commit.sha1);
      this.$router.push(this.$routes.dashboard.files());
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
        this.commits.forEach(commit => {
          commit.date = dateFormatter.default(commit.date);
          commit.message = truncate(commit.message, 15);
        });
      });
    }
  }
};
</script>
