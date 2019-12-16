<template>
  <div>
    <div v-if="parent && !parent.isDirectory">
      <pre>{{ childFile ? childFile.text : "" }}</pre>
    </div>
    <div v-else>
      <FileTable :files="children" @click-file-name="fileNameClicked"></FileTable>
    </div>
  </div>
</template>

<script>
import FileTable from "@/components/FileTable.vue";

let projectName = null;
let user = null;
export default {
  layout: "Project",
  components: {
    FileTable
  },
  data() {
    return {
      parent: null,
      childFile: null,
      children: []
    };
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

    console.log(this.$route.params);
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

    this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
      user = response.data;
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
        this.childFile = {
          text: response.data
        };
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
