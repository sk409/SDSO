<template>
  <div class="repository-files">
    <div
      v-for="(file, index) in files"
      :key="file.name"
      :class="{ 'border-bottom': index != files.length - 1 }"
      class="file-item p-2"
    >
      <i
        :class="{
          'el-icon-document': !file.isDirectory,
          'el-icon-folder': file.isDirectory
        }"
      ></i>
      <span class="file-item-name" @click="fileNameClicked(file)">{{
        file.name
      }}</span>
    </div>
  </div>
</template>

<script>
export default {
  name: "RepositoryFiles",
  props: {
    files: {
      type: Array,
      required: true
    }
  },
  methods: {
    fileNameClicked(file) {
      this.$emit("click-file-name", file);
    }
  }
};
</script>

<style scoped>
.repository-files {
  border: 1px solid #c0c0c0;
}

.file-item:hover {
  background: #f0f0f0;
}

.file-item-name {
  color: rgb(66, 120, 212);
}

.file-item-name:hover {
  cursor: pointer;
  text-decoration: underline;
}
</style>
