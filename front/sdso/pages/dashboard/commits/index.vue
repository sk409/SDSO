<template>
  <div>
    <v-subheader>コミット一覧</v-subheader>
    <v-divider class="mb-1"></v-divider>
    <GitToolbar
      class="mb-3"
      :hide-revision="true"
      @change-branchname="fetchCommits()"
    ></GitToolbar>
    <v-row justify="center" v-if="commits.length">
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
      commits: []
    };
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "git");
    this.fetchCommits();
  },
  methods: {
    ...mapMutations({
      setRevision: mutations.git.setRevision
    }),
    clickCommit(commit) {
      this.setRevision(commit.sha1);
      this.$router.push(this.$routes.dashboard.commits.show);
    },
    fetchCommits() {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
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
        teamname: team.name,
        projectname: project.name,
        branchname
      };
      ajax.get(url.base, data).then(response => {
        this.commits = response.data;
      });
    }
  }
};
</script>
