<template>
  <div ref="mainView" class="d-flex h-100">
    <div class="h-100 border-right">
      <slot name="sidemenu"></slot>
    </div>
    <div :style="contentStyle" class="w-100 h-100 overflow-x-hidden overflow-y-auto">
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
</style>
