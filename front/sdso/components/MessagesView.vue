<template>
  <div class="h-100 messages-view">
    <div ref="messages" class="messages overflow-y-auto" @scroll="scroll">
      <div v-if="more" class="text-center mt-3">
        <v-progress-circular indeterminate></v-progress-circular>
      </div>
      <div
        v-for="(message, index) in messages"
        :key="message.id"
        class="py-3 mx-auto message"
        @click="clickMessage(message, index)"
      >
        <div
          v-if="message.parents && message.parents.length != 0"
          class="pa-3 my-3 w-100 parents-card"
        >
          <div v-for="parent in message.parents" :key="parent.id">
            <MessageView :message="parent"></MessageView>
            <v-divider></v-divider>
          </div>
        </div>
        <MessageView :message="message"></MessageView>
        <MessageInput
          v-if="inputs[index].visible"
          :parent="message"
          :rows.sync="inputs[index].rows"
          :users="users"
          class="w-100 replay-message-input"
          @send="send"
        ></MessageInput>
        <v-divider class="mt-2"></v-divider>
      </div>
    </div>
    <MessageInput
      :rows.sync="rows"
      :users="users"
      class="message-input w-100"
      @send="send"
    ></MessageInput>
  </div>
</template>

<script>
import MessageInput from "@/components/MessageInput.vue";
import MessageView from "@/components/MessageView.vue";
import { count } from "@/assets/js/utils.js";
let fetchLength = 10;
export default {
  props: {
    loadMessages: {
      type: Function
    },
    messageCount: {
      type: Number,
      required: true
    },
    messages: {
      type: Array,
      default: []
    },
    postMessage: {
      type: Function,
      required: true
    },
    users: {
      type: Array
    }
  },
  components: {
    MessageInput,
    MessageView
  },
  data() {
    return {
      inputs: [],
      rows: 1
    };
  },
  computed: {
    more() {
      return (
        fetchLength < this.messageCount &&
        this.messages.length !== this.messageCount
      );
    }
  },
  mounted() {
    this.loadMessages(0, fetchLength, () => {
      this.$nextTick(() => {
        this.scrollToBottom();
      });
    });
  },
  watch: {
    messages(newMessages) {
      this.inputs = newMessages.map(message => ({
        rows: 1,
        visible: false
      }));
    }
  },
  methods: {
    clickMessage(message, index) {
      this.inputs.forEach(input => (input.visible = false));
      this.messages.forEach(message => {
        message.parents = [];
      });
      this.inputs[index].visible = true;
      const parents = [];
      let parent = message.parent
        ? this.messages.find(m => m.id === message.parent.id)
        : null;
      while (parent) {
        parents.push(parent);
        parent = parent.parent
          ? this.messages.find(m => m.id === parent.parent.id)
          : null;
      }
      message.parents = this.messages.filter(message =>
        parents.find(parent => parent.id === message.id)
      );
      if (index === this.messages.length - 1) {
        this.scrollToBottom();
      }
    },
    scroll() {
      if (this.$refs.messages.scrollTop === 0) {
        const preScrollHeight = this.$refs.messages.scrollHeight;
        this.loadMessages(
          this.messages.length,
          this.messages.length + fetchLength,
          () => {
            this.$nextTick(() => {
              this.$refs.messages.scrollTop +=
                this.$refs.messages.scrollHeight - preScrollHeight;
            });
          }
        );
      }
    },
    scrollToBottom() {
      const messages = this.$refs.messages;
      if (!messages) {
        return;
      }
      messages.scrollTop = messages.scrollHeight;
    },
    send(message, parent) {
      this.postMessage(message, parent, () => {
        this.$nextTick(() => {
          this.scrollToBottom();
        });
        this.$emit("update:messageCount", this.messageCount + 1);
      });
    }
  }
};
</script>

<style>
.message-input {
  position: absolute;
  bottom: 0;
  left: 0;
  border: 1px solid lightgray;
  height: 10%;
  max-height: 30%;
}
.message {
  width: 95%;
  position: relative;
}
.messages {
  height: 90%;
}
.messages-view {
  position: relative;
}
.parents-card {
  position: absolute;
  top: -10px;
  transform: translate(0, -100%);
  background: white;
  border: 1px solid lightgray;
  box-shadow: 0 0 8px gray;
}
.replay-message-input {
  max-height: 300px;
  border: 1px solid lightgray;
}
</style>
