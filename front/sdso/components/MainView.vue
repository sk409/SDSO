<template>
  <div ref="mainView" class="d-flex h-100">
    <div class="h-100 border-right sidemenu">
      <slot name="sidemenu"></slot>
    </div>
    <div
      :style="contentStyle"
      class="h-100 overflow-x-hidden overflow-y-auto content"
    >
      <slot name="content"></slot>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      contentStyle: {}
    };
  },
  mounted() {
    const maxHeight = this.$refs.mainView.clientHeight;
    this.contentStyle["max-height"] = maxHeight + "px";
  }
};
</script>

<style>
.border-right {
  border-right: 2px solid lightgrey;
}
.content {
  width: 80%;
}
.sidemenu {
  width: 20%;
}
</style>
