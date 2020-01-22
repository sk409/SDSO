<template>
  <div>
    <v-row justify="end">
      <v-col cols="3">
        <v-select
          v-model="branchname"
          :items="branchnames"
          label="ブランチを選択してください"
        ></v-select>
      </v-col>
      <v-col cols="3">
        <v-select
          v-model="commitSHA1"
          :items="commitSHA1s"
          label="コミットを選択してください"
        ></v-select>
      </v-col>
    </v-row>
    <v-divider></v-divider>
    <v-card>
      <v-card-text>
        <v-simple-table>
          <template v-slot:default>
            <tbody>
              <tr
                v-for="fileItem in fileItems"
                :key="fileItem.path"
                @click="clickFileItem(fileItem)"
              >
                <td>
                  <v-icon v-if="fileItem.isDirectory"
                    >mdi-folder-outline</v-icon
                  >
                  <v-icon v-else>mdi-file-document-box-outline</v-icon>
                  <span class="ml-3">{{ fileItem.name }}</span>
                </td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import { pathBranches, pathCommits, pathFiles, Url } from "@/assets/js/urls.js";
export default {
  layout: "dashboard",
  data() {
    return {
      branchname: "",
      branchnames: [],
      commits: [],
      commitSHA1: "",
      fileItems: [],
      user: null
    };
  },
  computed: {
    commitSHA1s() {
      return this.commits ? this.commits.map(commit => commit.sha1) : [];
    }
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "git");
    this.$store.subscribe((mutation, state) => {
      if (mutation.type !== mutations.projects.setProject) {
        return;
      }
      this.fetchData();
    });
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchData();
    });
  },
  methods: {
    clickFileItem(fileItem) {
      if (fileItem.isDirectory) {
        this.$router.push(
          this.$routes.dashboard.git.files(
            this.branchname,
            this.commitSHA1,
            fileItem.path
          )
        );
      }
    },
    fetchData() {
      const f = this.fetchBranches();
      if (!f) {
        return;
      }
      f.then(response => {
        this.branchnames = response.data;
        const branchname = this.$route.params.branchname;
        const commitSHA1 =
          this.$route.params.commitSHA1 || this.$route.params.pathMatch;
        if (
          !branchname ||
          !commitSHA1 ||
          !this.branchnames.includes(branchname)
        ) {
          return;
        }
        this.fetchCommits(branchname)
          .then(response => {
            const commits = response.data;
            if (!commits.some(commit => commit.sha1 === commitSHA1)) {
              return;
            }
            this.branchname = branchname;
            this.commits = commits;
            this.commitSHA1 = commitSHA1;
            const path = this.$route.params.commitSHA1
              ? this.$route.params.pathMatch
              : "";
            return this.fetchFiles(this.commitSHA1, path);
          })
          .then(response => {
            this.fileItems = response.data;
          });
      });
    },
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
      return ajax.get(url.base, data);
    },
    fetchCommits(branchname) {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathCommits);
      const data = {
        username: this.user.name,
        projectname: project.name,
        branchname
      };
      return ajax.get(url.base, data);
    },
    fetchFiles(treeIsh, path) {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      if (!path) {
        path = "";
      }
      const url = new Url(pathFiles);
      const data = {
        username: this.user.name,
        projectname: project.name,
        treeIsh: treeIsh,
        path: path
      };
      return ajax.get(url.base, data);
    }
  }
};
</script>
