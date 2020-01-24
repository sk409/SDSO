<template>
  <div>
    <v-subheader>コミット一覧</v-subheader>
    <v-divider class="mb-3"></v-divider>
    <GitToolbar
      class="mb-3"
      :hide-revision="true"
      @change-branchname="fetchCommits()"
    ></GitToolbar>
    <v-row justify="center">
      <v-col cols="11">
        <v-card class="mb-4">
          <v-simple-table>
            <template>
              <thead>
                <tr>
                  <th>SHA1</th>
                  <th>コミットメッセージ</th>
                  <th>日付</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="commit in commits"
                  :key="commit.sha1"
                  @click="clickCommit(commit)"
                >
                  <td>{{ commit.sha1 }}</td>
                  <td>{{ commit.message | truncate(15) }}</td>
                  <td>{{ commit.date | dateDefault }}</td>
                </tr>
              </tbody>
            </template>
          </v-simple-table>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import GitToolbar from "@/components/GitToolbar.vue";
import mutations from "@/assets/js/mutations.js";
import { pathCommits, Url } from "@/assets/js/urls.js";
import { dateFormatter, truncate } from "@/assets/js/utils.js";
import { mapMutations } from "vuex";
export default {
  layout: "dashboard",
  components: {
    GitToolbar
  },
  data() {
    return {
      commits: [],
      // tableHeaders: [
      //   { text: "SHA1", value: "sha1" },
      //   { text: "コミットメッセージ", value: "message" },
      //   { text: "日付", value: "date" }
      // ],
      user: null
    };
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "git");
    this.$fetchUser().then(response => {
      this.user = response.data;
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
      const branchname = this.$store.state.git.branchname;
      if (!branchname) {
        return;
      }
      const url = new Url(pathCommits);
      const data = {
        username: this.user.name,
        projectname: project.name,
        branchname
      };
      ajax.get(url.base, data).then(response => {
        this.commits = response.data;
        // this.commits.forEach(commit => {
        //   commit.date = dateFormatter.default(commit.date);
        //   commit.message = truncate(commit.message, 15);
        // });
      });
    }
  }
};
</script>
