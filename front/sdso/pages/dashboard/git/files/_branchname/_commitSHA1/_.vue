<template>
  <div>
    <v-row justify="end">
      <v-col cols="3">
        <v-select
          v-model="branchname"
          :items="branchnames"
          no-data-text="ブランチがありません"
          :label="branchnames.length + '個のブランチ'"
          @input="$router.push($routes.dashboard.git.files(branchname))"
        ></v-select>
      </v-col>
      <v-col cols="3">
        <v-select
          v-model="commitSHA1"
          :items="commitSHA1s"
          no-data-text="コミットがありません"
          @input="
            $router.push($routes.dashboard.git.files(branchname, commitSHA1))
          "
          :label="commitSHA1s.length + '個のコミット'"
        ></v-select>
      </v-col>
    </v-row>
    <v-card>
      <v-card-text>
        <pre v-if="fileMode" class="body-1">{{ fileText }}</pre>
        <v-simple-table v-else>
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
      fileMode: null,
      fileText: "",
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
      const route = this.$routes.dashboard.git.files(
        this.branchname,
        this.commitSHA1
      );
      if (route === this.$route.path) {
        this.fetchData();
      } else {
        this.$router.push(route);
      }
    });
  },
  mounted() {
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchData();
    });
  },
  methods: {
    clickFileItem(fileItem) {
      this.$router.push(
        this.$routes.dashboard.git.files(
          this.branchname,
          this.commitSHA1,
          fileItem.path,
          !fileItem.isDirectory
        )
      );
    },
    fetchData() {
      const f = this.fetchBranches();
      if (!f) {
        return;
      }
      f.then(response => {
        this.branchnames = response.data;
        this.branchname = this.$route.params.branchname;
        if (!this.branchname || !this.branchnames.includes(this.branchname)) {
          return;
        }
        this.fetchCommits(this.branchname).then(response => {
          const commits = response.data;
          this.commits = commits;
          this.commitSHA1 =
            this.$route.params.commitSHA1 || this.$route.params.pathMatch;
          if (
            !this.commitSHA1 ||
            !commits.some(commit => commit.sha1 === this.commitSHA1)
          ) {
            return;
          }
          const path = this.$route.params.commitSHA1
            ? this.$route.params.pathMatch
            : "";
          this.fileMode = this.$route.query.file === "true";
          if (this.fileMode) {
            this.fetchFileText(this.commitSHA1, path).then(response => {
              this.fileText = response.data;
            });
          } else {
            this.fetchFiles(this.commitSHA1, path).then(response => {
              const folders = response.data.filter(
                fileItem => fileItem.isDirectory
              );
              const files = response.data.filter(
                fileItem => !fileItem.isDirectory
              );
              this.fileItems = folders.concat(files);
            });
          }
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
    },
    fetchFileText(revision, path) {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathFiles);
      const data = {
        username: this.user.name,
        projectname: project.name,
        revision,
        path
      };
      return ajax.get(url.text, data);
    }
  }
};
</script>
