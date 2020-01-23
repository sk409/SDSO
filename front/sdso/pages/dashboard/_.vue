<template>
  <div>
    <v-row justify="center">
      <v-col cols="11">
        <v-card>
          <v-card-text>
            <pre v-if="file" class="body-1 code">{{ fileText }}</pre>
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
      </v-col>
    </v-row>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import { pathFiles, Url } from "@/assets/js/urls.js";
import { mapMutations } from "vuex";
export default {
  layout: "git",
  data() {
    return {
      completion: false,
      file: false,
      fileItems: [],
      fileText: "",
      user: null
    };
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "git");
    this.subscribe();
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.setup();
    });
  },
  methods: {
    ...mapMutations({
      setBranchname: mutations.git.setBranchname,
      setRevision: mutations.git.setRevision
    }),
    clickFileItem(fileItem) {
      this.$router.push(
        this.$routes.dashboard.files(fileItem.path, !fileItem.isDirectory)
      );
    },
    fetchFiles() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const revision = this.$store.state.git.revision;
      if (!revision) {
        return;
      }
      const path = this.$route.params.pathMatch;
      const url = new Url(pathFiles);
      const data = {
        username: this.user.name,
        projectname: project.name,
        treeIsh: revision,
        path: path
      };
      ajax.get(url.base, data).then(response => {
        const folders = response.data.filter(fileItem => fileItem.isDirectory);
        const files = response.data.filter(fileItem => !fileItem.isDirectory);
        this.fileItems = folders.concat(files);
        this.completion = true;
      });
    },
    fetchFileText() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const revision = this.$store.state.git.revision;
      if (!revision) {
        return;
      }
      const path = this.$route.params.pathMatch;
      const url = new Url(pathFiles);
      const data = {
        username: this.user.name,
        projectname: project.name,
        revision,
        path
      };
      ajax.get(url.text, data).then(response => {
        this.fileText = response.data;
        this.completion = true;
      });
    },
    subscribe() {
      this.$store.subscribe((mutation, state) => {
        switch (mutation.type) {
          case mutations.git.setBranchname:
            this.$router.push(this.$routes.dashboard.files());
            break;
          case mutations.git.setRevision:
            this.setup();
            break;
        }
      });
    },
    setup() {
      this.file = this.$route.query.file === "true";
      if (this.file) {
        this.fetchFileText();
      } else {
        this.fetchFiles();
      }
    }
  }
};
</script>

<style>
.code {
  width: 780px;
  overflow-x: scroll;
}
</style>
