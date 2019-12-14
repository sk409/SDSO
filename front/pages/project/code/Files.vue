<template>
  <div>
    <FileTable :files="files" @click-file-name="fileNameClicked"></FileTable>
  </div>
</template>

<script>
import FileTable from "@/components/FileTable.vue";

let user = null;
export default {
  layout: "Project",
  components: {
    FileTable
  },
  data() {
    return {
      files: []
    };
  },
  created() {
    this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
      user = response.data;
      this.fetchFiles("");
    });
  },
  methods: {
    fileNameClicked(file) {
      if (file.isDirectory) {
        this.fetchFiles(file.path);
        return;
      }
    },
    fetchFiles(path) {
      const data = {
        userName: user.Name,
        projectName: this.$route.params.name,
        path: path
      };
      this.$ajax.get(this.$urls.files, data, {}, response => {
        this.files = [];
        for (const fileName of Object.keys(response.data)) {
          this.files.push({
            path: response.data[fileName].path,
            name: fileName,
            isDirectory: response.data[fileName].isDirectory
          });
        }
        this.files = this.files.sort((a, b) => {
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
