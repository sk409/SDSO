<template>
  <div>
    <div class="mb-3">
      <!-- <div class="d-flex align-items-center">
        <div class="mr-1">ブランチ:</div>
        <el-select v-model="branch" @change="changeBranch">
          <el-option v-for="branch in branches" :key="branch" :label="branch" :value="branch"></el-option>
        </el-select>
      </div>-->
      <BranchSelection
        v-model="branchName"
        :branch-names="branchNames"
        @select-branch="changeBranch"
      ></BranchSelection>
    </div>
    <div v-if="parent && !parent.isDirectory">
      <div id="editor"></div>
    </div>
    <div v-else>
      <FileTable :files="children" @click-file-name="fileNameClicked"></FileTable>
    </div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.0/ace.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.0/ext-language_tools.js"></script>
  </div>
</template>

<script>
import BranchSelection from "@/components/BranchSelection.vue";
import FileTable from "@/components/FileTable.vue";
import { mapMutations } from "vuex";

let projectName = null;
let user = null;
let editor = null;
export default {
  name: "FileView",
  components: {
    BranchSelection,
    FileTable
  },
  data() {
    return {
      branchName: "",
      branchNames: [],
      childFile: null,
      children: [],
      parent: null
    };
  },
  computed: {
    pathParamUserName() {
      return this.$route.params.userName;
    },
    pathParamProjectName() {
      return this.$route.params.projectName
        ? this.$route.params.projectName
        : this.$route.params.pathMatch;
    },
    pathParamPath() {
      const pathPrefix =
        this.$routes.projectCode(
          this.pathParamUserName,
          this.pathParamProjectName
        ) + "/";
      const path = this.$route.path.startsWith(pathPrefix)
        ? this.$route.path.substring(pathPrefix.length)
        : "";
      return path;
    }
  },
  created() {
    this.branchName = this.$store.state.project.branchName;
    this.fetchData();
  },
  watch: {
    childFile(value) {
      const modes = {
        php: "php",
        js: "javascript"
      };
      const mode = Object.keys(modes).includes(value.extension)
        ? modes[value.extension]
        : "text";
      editor = ace.edit("editor");
      editor.$blockScrolling = Infinity;
      editor.setOptions({
        enableBasicAutocompletion: true,
        enableSnippets: false,
        enableLiveAutocompletion: true
      });
      editor.setTheme("ace/theme/monokai");
      editor.getSession().setMode("ace/mode/" + mode);
      editor.setFontSize(20);
      editor.setValue(this.childFile ? this.childFile.text : "");
    }
  },
  methods: {
    ...mapMutations({
      setBranchName: "project/setBranchName"
    }),
    changeBranch(newBranchName) {
      this.branchName = newBranchName;
      this.setBranchName(newBranchName);
      if (this.pathParamPath === "") {
        this.fetchData();
      } else {
        this.$router.push(
          this.$routes.projectCode(
            this.pathParamUserName,
            this.pathParamProjectName
          )
        );
      }
    },
    fileNameClicked(file) {
      this.$router.push(
        this.$routes.projectCode(
          this.pathParamUserName,
          this.pathParamProjectName,
          file.path
        )
      );
    },
    fetchData() {
      const userData = {
        name: this.pathParamUserName
      };
      this.$ajax.get(this.$urls.users, userData, {}, response => {
        if (response.status !== 200) {
          return;
        }
        user = response.data[0];
        if (this.pathParamPath === "") {
          this.fetchFiles(this.pathParamPath);
        } else {
          let parentPath = "";
          if (this.pathParamPath.includes("/")) {
            const pathComponents = this.pathParamPath.split("/");
            parentPath = pathComponents
              .slice(0, pathComponents.length - 1)
              .join("/");
          }
          const data = {
            userName: this.pathParamUserName,
            projectName: this.pathParamProjectName,
            branchName: this.$store.state.project.branchName,
            path: parentPath
          };
          this.$ajax.get(this.$urls.files, data, {}, response => {
            for (const fileName of Object.keys(response.data)) {
              if (response.data[fileName].path === this.pathParamPath) {
                this.parent = response.data[fileName];
                break;
              }
            }
            if (this.parent && this.parent.isDirectory) {
              this.fetchFiles(this.pathParamPath);
            } else {
              this.fetchFileText(this.pathParamPath);
            }
          });
        }
      });
      const branchData = {
        userName: this.pathParamUserName,
        projectName: this.pathParamProjectName
      };
      this.$ajax.get(this.$urls.branches, branchData, {}, response => {
        if (response.status !== 200) {
          return;
        }
        this.branchNames = response.data;
      });
    },
    fetchFileText(path) {
      const data = {
        userName: this.pathParamUserName,
        projectName: this.pathParamProjectName,
        branchName: this.$store.state.project.branchName,
        path: path
      };
      this.$ajax.get(this.$urls.filesText, data, {}, response => {
        const s = this.$route.path.split(".");
        this.childFile = {
          text: response.data
        };
        if (s.length) {
          this.childFile.extension = s[s.length - 1];
        }
      });
    },
    fetchFiles(path) {
      const data = {
        userName: this.pathParamUserName,
        projectName: this.pathParamProjectName,
        branchName: this.$store.state.project.branchName,
        path: path
      };
      this.$ajax.get(this.$urls.files, data, {}, response => {
        this.children = [];
        for (const fileName of Object.keys(response.data)) {
          this.children.push({
            path: response.data[fileName].path,
            name: fileName,
            isDirectory: response.data[fileName].isDirectory
          });
        }
        this.children = this.children.sort((a, b) => {
          if (a.isDirectory && !b.isDirectory) {
            return -1;
          }
          if (!a.isDirectory && b.isDirectory) {
            return 1;
          }
          return a.name < b.name ? -1 : 1;
        });
      });
    }
  }
};
</script>

<style scoped>
#editor {
  height: 500px;
}
</style>
