<template>
  <div style="position:relative;">
    <v-card v-show="usersCard && mentionUsers.length" class="users-card w-100">
      <div v-for="user in mentionUsers" :key="user.id" class="user" @click="clickUser(user)">
        <v-avatar>
          <v-img :src="$serverUrl(user.profileImagePath)"></v-img>
        </v-avatar>
        <span>{{user.name}}</span>
      </div>
    </v-card>
    <textarea
      ref="textarea"
      v-model="message"
      :rows="rows"
      class="resize-none w-100 h-100"
      placeholder="Enterで改行してShift+Enterで送信します。"
      @input="input"
      @keydown.delete="shrink($event)"
      @keydown.enter.exact="grow($event)"
      @keydown.enter.shift="send"
    ></textarea>
  </div>
</template>

<script>
import { count } from "@/assets/js/utils.js";
export default {
  props: {
    parent: {
      type: Object,
      default: null
    },
    rows: {
      type: Number,
      default: 1
    },
    users: {
      type: Array
    }
  },
  data() {
    return {
      message: "",
      search: "",
      usersCard: false
    };
  },
  computed: {
    mentionUsers() {
      return this.users.filter(user => user.name.startsWith(this.search));
    }
  },
  methods: {
    clickUser(user) {
      const textarea = this.$refs.textarea;
      let message =
        this.message.substr(0, textarea.selectionStart - this.search.length) +
        user.name;
      if (textarea.selectionStart < this.message.length) {
        message += this.message.substr(textarea.selectionStart);
      }
      this.message = message;
      this.usersCard = false;
    },
    grow(e) {
      if (e.isComposing) {
        return;
      }
      this.$emit("update:rows", this.rows + 1);
      const textarea = this.$refs.textarea;
      if (textarea.selectionStart !== textarea.selectionEnd) {
        const deletedText = this.message.substring(
          textarea.selectionStart,
          textarea.selectionEnd
        );
        this.$emit("update:rows", this.rows - count(deletedText, "\n"));
      }
    },
    input(e) {
      this.search = "";
      const textarea = this.$refs.textarea;
      let index = textarea.selectionStart - 1;
      while (true) {
        if (index < 0) {
          this.usersCard = false;
          break;
        }
        const character = this.message[index];
        if (character === "@") {
          this.usersCard = true;
          break;
        }
        const regex = /[a-zA-Z0-9]/;
        if (character === "\n" || !regex.test(character)) {
          this.usersCard = false;
          break;
        }
        this.search = character + this.search;
        index -= 1;
      }
    },
    send(e) {
      e.preventDefault();
      this.$emit("send", this.message, this.parent);
      this.$emit("update:rows", 1);
      this.message = "";
    },
    shrink(e) {
      if (e.isComposing) {
        return;
      }
      const textarea = this.$refs.textarea;
      if (textarea.selectionStart === textarea.selectionEnd) {
        const char = this.message[textarea.selectionStart - 1];
        if (char === "\n") {
          this.$emit("update:rows", this.rows - 1);
        }
      } else {
        const deletedText = this.message.substring(
          textarea.selectionStart,
          textarea.selectionEnd
        );
        this.$emit("update:rows", this.rows - count(deletedText, "\n"));
      }
    }
  }
};
</script>

<style>
.user {
  cursor: pointer;
}
.user:hover {
  background: rgb(240, 240, 240);
}
.users-card {
  position: absolute;
  left: 0;
  top: 0;
  transform: translate(0, -100%);
}
</style>