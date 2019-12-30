<template>
  <div>
    <div v-if="parent && !parent.isDirectory">
      <div id="editor"></div>
    </div>
    <div v-else>
      <FileTable
        :files="children"
        @click-file-name="fileNameClicked"
      ></FileTable>
    </div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.0/ace.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.0/ext-language_tools.js"></script>
  </div>
</template>

<script>
import FileTable from "@/components/FileTable.vue";

let projectName = null;
let user = null;
let editor = null;
export default {
  layout: "Project",
  components: {
    FileTable
  },
  data() {
    return {
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
    }
  },
  created() {
    // this.$ajax.get("test", {}, {}, response => {
    //   for (const fileName of Object.keys(response.data)) {
    //     if (response.data[fileName].path === path) {
    //       this.parent = response.data[fileName];
    //       break;
    //     }
    //   }
    //   if (this.parent && this.parent.isDirectory) {
    //     this.fetchFiles(path);
    //   } else {
    //     this.fetchFileText(path);
    //   }
    // });

    // console.log(this.$route.params);
    projectName = this.$route.params.projectName
      ? this.$route.params.projectName
      : this.$route.params.pathMatch;
    // this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
    //   user = response.data;
    //   const pathPrefix = this.$routes.projectCode(projectName) + "/";
    //   const path = this.$route.path.startsWith(pathPrefix)
    //     ? this.$route.path.substring(pathPrefix.length)
    //     : "";
    //   this.fetchFiles(path);
    // });

    const data = {
      name: this.pathParamUserName
    };
    this.$ajax.get(this.$urls.users, data, {}, response => {
      if (response.status !== 200) {
        return;
      }
      user = response.data[0];
      const pathPrefix = this.$routes.projectCode(user.Name, projectName) + "/";
      const path = this.$route.path.startsWith(pathPrefix)
        ? this.$route.path.substring(pathPrefix.length)
        : "";
      if (path === "") {
        this.fetchFiles(path);
      } else {
        let parentPath = "";
        if (path.includes("/")) {
          const pathComponents = path.split("/");
          parentPath = pathComponents
            .slice(0, pathComponents.length - 1)
            .join("/");
        }
        const data = {
          userName: user.Name,
          projectName: projectName,
          path: parentPath
        };
        this.$ajax.get(this.$urls.files, data, {}, response => {
          for (const fileName of Object.keys(response.data)) {
            if (response.data[fileName].path === path) {
              this.parent = response.data[fileName];
              break;
            }
          }
          if (this.parent && this.parent.isDirectory) {
            this.fetchFiles(path);
          } else {
            this.fetchFileText(path);
          }
        });
      }
    });
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
    fileNameClicked(file) {
      this.$router.push(
        this.$routes.projectCode(user.Name, projectName, file.path)
      );
    },
    fetchFileText(path) {
      const data = {
        userName: user.Name,
        projectName: projectName,
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
        userName: user.Name,
        projectName: projectName,
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
      // const data = {
      //   userName: user.Name,
      //   projectName: projectName,
      //   path: path
      // };
      // this.$ajax.get(this.$urls.files, data, {}, response => {
      // this.children = [];
      // for (const fileName of Object.keys(response.data)) {
      //   this.children.push({
      //     path: response.data[fileName].path,
      //     name: fileName,
      //     isDirectory: response.data[fileName].isDirectory
      //   });
      // }
      // this.children = this.children.sort((a, b) => {
      //   if (a.isDirectory && !b.isDirectory) {
      //     return -1;
      //   }
      //   if (!a.isDirectory && b.isDirectory) {
      //     return 1;
      //   }
      //   return a.name < b.name ? -1 : 1;
      // });
      // });
    }
  }
};
</script>

<style scoped>
#editor {
  height: 500px;
}
</style>
