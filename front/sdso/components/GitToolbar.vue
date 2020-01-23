<template>
  <v-container fluid>
    <v-row justify="end">
      <v-col cols="3">
        <v-select
          :items="branchnames"
          no-data-text="ブランチがありません"
          :value="branchname"
          :label="branchnames.length + '個のブランチ'"
          @input="selectBranchname"
        ></v-select>
      </v-col>
      <v-col cols="3">
        <v-select
          :items="revisions"
          no-data-text="コミットがありません"
          :value="revision"
          @input="selectRevision"
          :label="revisions.length + '個のコミット'"
        ></v-select>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  props: {
    branchname: {
      type: String,
      default: ""
    },
    branchnames: {
      type: Array,
      default: () => []
    },
    revision: {
      type: String,
      default: ""
    },
    revisions: {
      type: Array,
      default: () => []
    }
  },
  methods: {
    selectBranchname(branchname) {
      this.$emit("update:branchname", branchname);
      this.$emit("select-branchname", branchname);
    },
    selectRevision(revision) {
      this.$emit("update:revision", revision);
      this.$emit("select-revision", revision);
    }
  }
};
</script>
