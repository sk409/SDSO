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
          <MessagesView
            :load-messages="fetchMessages"
            :message-count.sync="messageCount"
            :messages="messages"
            :post-message="storeMessage"
            :users="users"
          ></MessagesView>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import MessagesView from "@/components/MessagesView.vue";
import {
  pathProjects,
  pathTests,
  pathTestMessages,
  Url
} from "@/assets/js/urls.js";
import { setupTest } from "@/assets/js/utils.js";

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
      messageCount: 0,
      messages: [],
      tabActive: "",
      tabs: ["コマンド", "コメント"],
      test: null,
      users: []
    };
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
      this.fetchMessageCount();
      this.fetchTestAndUsers();
      this.setupSocketMessage();
      this.setupSocketTest();
    });
  },
  methods: {
    fetchMessageCount() {
      const testId = this.$route.params.id;
      const url = new Url(pathTestMessages);
      const data = {
        testId
      };
      ajax.get(url.count, data).then(response => {
        this.messageCount = Number(response.data);
      });
    },
    fetchMessages(start, end, completion) {
      const testId = this.$route.params.id;
      const url = new Url(pathTestMessages);
      const data = {
        start,
        end,
        testId
      };
      ajax.get(url.range, data).then(response => {
        this.messages = response.data.concat(this.messages);
        completion();
      });
    },
    fetchTestAndUsers() {
      const url = new Url(pathTests);
      const data = {
        id: this.$route.params.id
      };
      ajax.get(url.base, data).then(response => {
        this.test = response.data[0];
        setupTest(this.test);
        const url = new Url(pathProjects);
        const data = {
          id: this.test.projectId
        };
        ajax.get(url.base, data).then(response => {
          this.users = response.data[0].users;
        });
      });
    },
    setupSocketMessage() {
      const that = this;
      const url = new Url(pathTestMessages);
      socket = new WebSocket(url.socket(user.id));
      socket.onmessage = function(e) {
        if (!that.test || !user) {
          return;
        }
        const message = JSON.parse(e.data);
        if (message.testId !== that.test.id) {
          return;
        }
        if (message.userId === user.id) {
          return;
        }
        that.messages.push(message);
      };
    },
    setupSocketTest() {
      const that = this;
      const url = new Url(pathTests);
      socket = new WebSocket(url.socket(user.id));
      socket.onmessage = function(e) {
        const test = JSON.parse(e.data);
        if (that.test.id !== test.id) {
          return;
        }
        setupTest(test);
        that.test = test;
      };
    },
    storeMessage(message, parent, completion) {
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
        completion();
      });
    }
  }
};
</script>

<style>
pre {
  font-family: "Roboto", sans-serif;
}
.messages-view {
  height: 450px;
}
</style>
