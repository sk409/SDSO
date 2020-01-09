<template>
  <div class="container-fluid">
    <!-- <div class="row">
      <div
        class="d-flex align-items-center justify-content-between border mt-3 mb-3 p-2 col-12 col-md-6 offset-md-3"
      >
        <div class="link mr-3" @click="transitionToFilesScene">ファイル</div>
        <div class="link mr-3" @click="transitionToCommitsScene">コミット</div>
        <div class="link" @click="transitionToBranchesScene">ブランチ</div>
      </div>
    </div>-->
    <div class="row mt-3">
      <el-tabs v-model="activeTabName" type="card" class="col">
        <el-tab-pane label="ファイル" name="files">
          <FileView></FileView>
        </el-tab-pane>
        <el-tab-pane label="コミット" name="commit">
          <CommitView></CommitView>
        </el-tab-pane>
        <el-tab-pane label="ブランチ" name="branch">
          <BranchView></BranchView>
        </el-tab-pane>
      </el-tabs>
      <!-- <div v-if="isFilesScene" class="col">
        <FileView></FileView>
      </div>
      <div v-else-if="isCommitsScene" class="col">
        <CommitView></CommitView>
      </div>
      <div v-else class="col">
        <BranchView></BranchView>
      </div>-->
    </div>
  </div>
</template>

<script>
import BranchView from "@/components/BranchView.vue";
import CommitView from "@/components/CommitView.vue";
import FileView from "@/components/FileView.vue";

export default {
  layout: "Project",
  components: {
    BranchView,
    CommitView,
    FileView
  },
  data() {
    return {
      activeTabName: "files"
    };
  }
  // computed: {
  //   isFilesScene() {
  //     return this.scene === sceneFiles;
  //   },
  //   isCommitsScene() {
  //     return this.scene === sceneCommits;
  //   },
  //   isBranchesScene() {
  //     return this.scene === sceneBranches;
  //   }
  // }
};
</script>
