<template>
  <div class="h-100 messages-view">
    <div ref="messagesContainer" class="messages overflow-y-auto" @scroll="scroll">
      <div v-if="more" class="text-center">
        <v-progress-circular indeterminate></v-progress-circular>
      </div>
      <div
        ref="messages"
        v-for="(message, index) in messages"
        :key="message.id"
        :message-id="message.id"
        :class="{'new-message': message.new}"
        class="message-container"
        @click="clickMessage(message, index)"
      >
        <div
          v-if="message.newMessageChip"
          :style="newMessageChipStyle"
          class="new-message-chip d-flex align-center"
        >
          <span>新着メッセージ</span>
        </div>
        <div class="pt-3 mx-auto message">
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
    </div>
    <MessageInput :rows.sync="rows" :users="users" class="message-input w-100" @send="send"></MessageInput>
  </div>
</template>

<script>
import MessageInput from "@/components/MessageInput.vue";
import MessageView from "@/components/MessageView.vue";
import { count } from "@/assets/js/utils.js";
let fetchLength = 10;
export default {
  props: {
    baselineMessage: {
      validator: v => typeof v === "object" || v === null
    },
    loadMessageIds: {
      type: Function
    },
    loadMessages: {
      type: Function
    },
    messageIds: {
      type: Array
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
      newMessageChipHeight: 48,
      rows: 1
    };
  },
  computed: {
    more() {
      return (
        fetchLength < this.messageIds.length &&
        this.messages.length !== this.messageIds.length
      );
    },
    newMessageChipStyle() {
      return {
        height: this.newMessageChipHeight + "px"
      };
    }
  },
  mounted() {
    if (this.loadMessageIds) {
      this.loadMessageIds(() => {
        this.$nextTick(() => {
          if (this.loadMessages) {
            if (this.baselineMessage) {
              const messageIdIndex = this.messageIds.findIndex(
                messageId => messageId === this.baselineMessage.id
              );
              const upperBoundMessageIdIndex = Math.min(
                messageIdIndex + fetchLength,
                this.messageIds.length
              );
              const startMessageIdIndex =
                messageIdIndex -
                (fetchLength - (upperBoundMessageIdIndex - messageIdIndex));
              const ids = this.messageIds.slice(
                Math.max(startMessageIdIndex - 1, 0),
                this.messageIds.length
              );
              this.loadMessages(ids, () => {
                this.$nextTick(() => {
                  const baselineMessageElement = this.$refs.messages.find(
                    message => {
                      const messageId = Array.from(message.attributes).find(
                        attribute => attribute.name === "message-id"
                      ).value;
                      return this.baselineMessage.id == messageId;
                    }
                  );
                  this.$refs.messagesContainer.scrollTop =
                    baselineMessageElement.getBoundingClientRect().top -
                    this.$refs.messagesContainer.getBoundingClientRect().top -
                    this.newMessageChipHeight / 2;
                });
              });
            } else {
              const ids = this.messageIds.slice(
                Math.max(this.messageIds.length - fetchLength, 0),
                this.messageIds.length
              );
              this.loadMessages(ids, () => {
                this.$nextTick(() => {
                  this.scrollToBottom();
                });
              });
            }
          }
        });
      });
    }
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
      if (this.$refs.messagesContainer.scrollTop === 0) {
        const preScrollHeight = this.$refs.messagesContainer.scrollHeight;
        const ids = this.messageIds.slice(
          Math.max(
            this.messageIds.length - this.messages.length - fetchLength,
            0
          ),
          this.messageIds.length - this.messages.length
        );
        this.loadMessages(ids, () => {
          this.$nextTick(() => {
            this.$refs.messagesContainer.scrollTop +=
              this.$refs.messagesContainer.scrollHeight - preScrollHeight;
          });
        });
      }
    },
    scrollToBottom() {
      const messagesContainer = this.$refs.messagesContainer;
      if (!messagesContainer) {
        return;
      }
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    },
    send(message, parent) {
      this.postMessage(message, parent, () => {
        this.$nextTick(() => {
          this.scrollToBottom();
        });
        // this.$emit("update:messageCount", this.messageCount + 1);
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
.message-container {
  position: relative;
}
.messages {
  max-height: 90%;
}
.messages-view {
  position: relative;
}
.new-message {
  background: rgb(255, 253, 228);
}
.new-message-chip {
  background: white;
  position: absolute;
  top: 0;
  left: 50%;
  transform: translate(-50%, -50%);
  box-shadow: 0 2px 3px 1px rgba(0, 0, 0, 0.5);
  padding: 0 1rem;
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
