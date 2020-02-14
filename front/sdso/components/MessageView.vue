<template>
  <div class="message">
    <div class="d-flex align-center">
      <v-avatar size="64">
        <v-img :src="$serverUrl(message.user.profileImagePath)"></v-img>
      </v-avatar>
      <span class="ml-3 title">{{ message.user.name }}</span>
      <span class="ml-auto mr-3 body-1">{{ time }}</span>
    </div>
    <pre class="my-3">{{ message.text }}</pre>
  </div>
</template>

<script>
import { dateFormatter } from "@/assets/js/utils.js";
let timer = null;
export default {
  props: {
    message: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      time: ""
    };
  },
  watch: {
    message: {
      immediate: true,
      handler(message) {
        this.time = dateFormatter.ago(message.createdAt);
        clearInterval(timer);
        setInterval(() => {
          this.time = dateFormatter.ago(message.createdAt);
        }, 1000);
      }
    }
  }
};
</script>

<style>
.message {
  position: relative;
}
</style>