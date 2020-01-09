<template>
  <div>
    <div
      class="d-flex align-items-center justify-content-center border py-2 selected-branch-container"
      @click="toggleSelectionDialog"
    >
      <div class="mr-2">ブランチ:</div>
      <div class="selected-branch">{{value}}</div>
    </div>
    <div v-show="selectionDialog.visible" class="p-2 shadow selection-dialog">
      <div class="text-center h6">ブランチ選択</div>
      <el-divider class="my-2"></el-divider>
      <div>
        <el-input v-model="inputBranchName"></el-input>
      </div>
      <el-divider class="my-2"></el-divider>
      <div
        v-for="branchName in filteredBranchNames"
        :key="branchName"
        class="p-2 branch"
        @click="selectBranch(branchName)"
      >{{branchName}}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: "BranchSelection",
  props: {
    value: {
      type: String,
      default: ""
    },
    branchNames: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      inputBranchName: "",
      selectionDialog: {
        visible: false
      }
    };
  },
  computed: {
    filteredBranchNames() {
      if (this.inputBranchName === "") {
        return this.branchNames;
      }
      return this.branchNames.filter(branchName =>
        branchName.includes(this.inputBranchName)
      );
    }
  },
  methods: {
    selectBranch(selectedBranchName) {
      this.$emit("select-branch", selectedBranchName);
      this.selectionDialog.visible = false;
    },
    toggleSelectionDialog() {
      this.selectionDialog.visible = !this.selectionDialog.visible;
    }
  }
};
</script>

<style scoped>
.selected-branch-container {
  width: 160px;
  background: rgb(240, 240, 240);
  cursor: pointer;
}

.select-box:hover {
  background: rgb(200, 200, 200);
}

.selected {
  font-weight: bold;
}

.selection-dialog {
  position: absolute;
  z-index: 100;
  transform: translateY(5px);
  background: rgb(240, 240, 240);
  width: 200px;
  border: 1px solid rgb(180, 180, 180);
}

.branch {
}

.branch:hover {
  color: white;
  background: rgb(42, 103, 207);
}
</style>