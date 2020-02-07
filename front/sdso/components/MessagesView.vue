<template>
  <div class="h-100 messages-view">
    <div ref="messages" class="messages overflow-y-auto">
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
          class="w-100 replay-message-input"
          @send="send"
        ></MessageInput>
        <v-divider></v-divider>
      </div>
    </div>
    <MessageInput
      :rows.sync="rows"
      class="message-input w-100"
      @send="send"
    ></MessageInput>
  </div>
</template>

<script>
import MessageInput from "@/components/MessageInput.vue";
import MessageView from "@/components/MessageView.vue";
import { count } from "@/assets/js/utils.js";
export default {
  props: {
    messages: {
      type: Array,
      required: true
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
  mounted() {
    this.$nextTick(() => {
      this.scroll();
    });
  },
  watch: {
    messages: {
      immediate: true,
      handler(messages) {
        this.inputs = messages.map(message => ({ rows: 1, visible: false }));
        this.$nextTick(() => {
          this.scroll();
        });
      }
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
      let parent = message.parent;
      while (parent) {
        parents.push(parent);
        parent = parent.parent;
      }
      message.parents = parents.reverse();
    },
    scroll() {
      const messages = this.$refs.messages;
      if (!messages) {
        return;
      }
      messages.scrollTop = messages.scrollHeight;
    },
    send(message, parent) {
      this.$emit("send", message, parent);
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
