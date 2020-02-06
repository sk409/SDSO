<template>
  <div>
    <div class="d-flex align-center">
      <v-avatar size="64">
        <v-img :src="$serverUrl(message.user.profileImagePath)"></v-img>
      </v-avatar>
      <span class="ml-3 title">{{ message.user.name }}</span>
      <span class="ml-auto mr-3 body-1">{{ message.createdAt | dateAgo }}</span>
    </div>
    <pre class="my-3">{{ message.text }}</pre>
  </div>
</template>

<script>
export default {
  props: {
    message: {
      type: Object,
      required: true
    }
  }
};
</script>
