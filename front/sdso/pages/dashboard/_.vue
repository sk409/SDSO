<template>
  <MainView>
    <template v-slot:sidemenu>
      <v-list class="pa-2">
        <v-list-item
          v-for="sidemenuItem in sidemenuItems"
          :key="sidemenuItem.title"
          :to="sidemenuItem.route"
          router
        >
          <v-list-item-action>
            <v-icon>{{ sidemenuItem.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ sidemenuItem.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </template>
    <template v-slot:content>
      <div class="h-100">
        <v-subheader>ファイル一覧</v-subheader>
        <v-divider class="mb-1"></v-divider>
        <GitToolbar
          class="mb-3"
          @change-branchname="changeBranchname"
          @change-revision="changeRevision"
        ></GitToolbar>
        <div v-if="file" id="editor"></div>
        <v-row v-else-if="completion" justify="center">
          <v-col cols="11">
            <v-card class="mb-4">
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
          </v-col>
        </v-row>
      </div>
    </template>
  </MainView>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import GitToolbar from "@/components/GitToolbar.vue";
import MainView from "@/components/MainView.vue";
import mutations from "@/assets/js/mutations.js";
import { pathFiles, Url } from "@/assets/js/urls.js";
import { mapMutations } from "vuex";

const extension = str => {
  const components = str.split(".");
  const extension =
    components.length == 0 ? "" : components[components.length - 1];
  return extension;
};

const aceMode = path => {
  const ext = extension(path);
  const modes = {
    go: "golang",
    js: "javascript",
    php: "php"
  };
  const base = "ace/mode/";
  if (!modes[ext]) {
    return base + "text";
  }
  return base + modes[ext];
};

let editor = null;
export default {
  layout: "dashboard",
  components: {
    GitToolbar,
    MainView
  },
  data() {
    return {
      completion: false,
      file: false,
      fileItems: [],
      fileText: "",
      sidemenuItems: [
        {
          title: "コミット",
          icon: "mdi-source-commit",
          route: this.$routes.dashboard.commits.base
        },
        {
          title: "ファイル",
          icon: "mdi-file-outline",
          route: this.$routes.dashboard.files()
        }
      ],
      user: null
    };
  },
  created() {
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.setup();
    });
  },
  methods: {
    changeBranchname() {
      this.$router.push(this.$routes.dashboard.files());
    },
    changeRevision() {
      this.setup();
    },
    clickFileItem(fileItem) {
      this.$router.push(
        this.$routes.dashboard.files(fileItem.path, !fileItem.isDirectory)
      );
    },
    fetchFiles() {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
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
        teamname: team.name,
        projectname: project.name,
        revision,
        path
      };
      ajax.get(url.base, data).then(response => {
        const folders = response.data.filter(fileItem => fileItem.isDirectory);
        const files = response.data.filter(fileItem => !fileItem.isDirectory);
        this.fileItems = folders.concat(files);
        this.completion = true;
      });
    },
    fetchFileText() {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
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
        teamname: team.name,
        projectname: project.name,
        revision,
        path
      };
      ajax.get(url.text, data).then(response => {
        this.fileText = response.data;
        editor = ace.edit("editor");
        editor.$blockScrolling = Infinity;
        editor.setTheme("ace/theme/xcode");
        editor.setFontSize(20);
        editor.setReadOnly(true);
        editor.setValue(response.data);
        editor.getSession().setMode(aceMode(path));
        this.completion = true;
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
#editor {
  width: 100%;
  height: 85%;
}
</style>
