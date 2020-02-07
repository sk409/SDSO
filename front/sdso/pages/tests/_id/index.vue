<template>
  <v-container v-if="test" class="mt-4">
    <div class="title">テストの詳細</div>
    <v-divider class="mb-5"></v-divider>
    <div class="mb-5">
      <v-chip :color="test.color" text-color="white" class="title">
        {{ test.status.toUpperCase() }}
      </v-chip>
    </div>
    <v-card class="pb-5">
      <v-tabs v-model="tabActive" class="mb-3">
        <v-tab v-for="tab in tabs" :key="tab" :href="`#${tab}`">
          {{ tab }}
        </v-tab>
      </v-tabs>
      <v-tabs-items v-model="tabActive" class="mx-5">
        <v-tab-item value="コマンド">
          <v-expansion-panels flat multiple :accordian="false">
            <v-expansion-panel
              v-for="result in test.results"
              :key="result.id"
              :style="{ 'border-left': `5px solid ${result.color}` }"
              class="my-2"
              style="border: 1px solid lightgray;"
            >
              <v-expansion-panel-header class="body-1">
                {{ result.command }}
              </v-expansion-panel-header>
              <v-expansion-panel-content>
                <pre
                  class="pa-3 blue-grey darken-4 white--text overflow-x-auto"
                  >{{ result.output }}</pre
                >
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-tab-item>
        <v-tab-item value="コメント">
          <MessagesView :messages="messages" @send="sendMessage"></MessagesView>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import MessagesView from "@/components/MessagesView.vue";
import { pathTests, pathTestMessages, Url } from "@/assets/js/urls.js";

let socket = null;
let user = null;
const statusLoading = "statusLoading";
const statusPermissionError = "statusPermissionError";
const statusOK = "statusOK";
export default {
  layout: "auth",
  components: {
    MessagesView
  },
  data() {
    return {
      messages: [],
      tabActive: "",
      tabs: ["コマンド", "コメント"],
      test: null
    };
  },
  created() {
    this.setupSocket();
    this.fetchMessages();
    this.$fetchUser().then(response => {
      user = response.data;
      const url = new Url(pathTests);
      const data = {
        id: this.$route.params.id
      };
      ajax.get(url.base, data).then(response => {
        this.test = response.data[0];
      });
    });
  },
  methods: {
    fetchMessages() {
      const testId = this.$route.params.id;
      const url = new Url(pathTestMessages);
      const data = {
        testId
      };
      ajax.get(url.base, data).then(response => {
        this.messages = response.data;
      });
    },
    sendMessage(message, parent) {
      const testId = this.$route.params.id;
      const url = new Url(pathTestMessages);
      const data = {
        text: message,
        testId,
        userId: user.id
      };
      if (parent) {
        data.parentId = parent.id;
      }
      ajax.post(url.base, data).then(response => {
        this.messages.push(response.data);
      });
    },
    setupSocket() {
      if (!WebSocket) {
        alert("WebSocketに対応していないブラウザです。");
        return;
      }
      const that = this;
      const url = new Url(pathTests);
      socket = new WebSocket(url.socket);
      socket.onmessage = function(e) {
        const test = JSON.parse(e.data);
        if (that.test.id !== test.id) {
          return;
        }
        that.test = test;
      };
    }
  }
};
</script>

<style>
.messages-view {
  height: 450px;
}
</style>
