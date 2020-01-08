<template>
  <div class="container-fluid">
    <div class="row">
      <div
        class="d-flex align-items-center justify-content-between border mt-3 mb-3 p-2 col-12 col-md-6 offset-md-3"
      >
        <div class="link mr-3" @click="transitionToFilesScene">ファイル</div>
        <div class="link mr-3" @click="transitionToCommitsScene">コミット</div>
        <div class="link" @click="transitionToBranchesScene">ブランチ</div>
      </div>
    </div>
    <div class="row mt-3">
      <div v-if="isFilesScene" class="col">
        <FileView></FileView>
      </div>
      <div v-else-if="isCommitsScene" class="col">
        <CommitView></CommitView>
      </div>
      <div v-else class="col">
        <BranchView></BranchView>
      </div>
    </div>
  </div>
</template>

<script>
import BranchView from "@/components/BranchView.vue";
import CommitView from "@/components/CommitView.vue";
import FileView from "@/components/FileView.vue";

const sceneFiles = "sceneFiles";
const sceneCommits = "sceneCommits";
const sceneBranches = "sceneBranches";
export default {
  layout: "Project",
  components: {
    BranchView,
    CommitView,
    FileView
  },
  data() {
    return {
      scene: sceneFiles
    };
  },
  computed: {
    isFilesScene() {
      return this.scene === sceneFiles;
    },
    isCommitsScene() {
      return this.scene === sceneCommits;
    },
    isBranchesScene() {
      return this.scene === sceneBranches;
    }
  },
  methods: {
    transitionToFilesScene() {
      this.scene = sceneFiles;
    },
    transitionToCommitsScene() {
      this.scene = sceneCommits;
    },
    transitionToBranchesScene() {
      this.scene = sceneBranches;
    }
  }
};
</script>
