<template>
  <v-snackbar v-model="snackbar" :timeout="2000" top>
    <span>{{ $store.state.notifications.message }}</span>
    <v-btn icon @click="close">
      <v-icon>mdi-close</v-icon>
    </v-btn>
  </v-snackbar>
</template>

<script>
import mutations from "@/assets/js/mutations.js";
import { mapMutations } from "vuex";
let unsubscribe = null;
export default {
  data() {
    return {
      snackbar: false
    };
  },
  created() {
    this.snackbar = this.$store.state.notifications.message !== "";
    unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === mutations.notifications.setMessage) {
        this.snackbar = state.notifications.message !== "";
      }
    });
  },
  destroyed() {
    unsubscribe();
  },
  methods: {
    ...mapMutations({
      setMessage: mutations.notifications.setMessage
    }),
    close() {
      this.snackbar = false;
      this.setMessage("");
    }
  }
};
</script>
