<template>
  <v-container v-if="vulnerability" class="mt-4">
    <div class="title">脆弱性の詳細</div>
    <v-divider class="mb-7"></v-divider>
    <v-card class="pb-5">
      <v-tabs v-model="tabActive" class="mb-3">
        <v-tab v-for="tab in tabs" :key="tab" :href="`#${tab}`">{{
          tab
        }}</v-tab>
      </v-tabs>
      <v-tabs-items v-model="tabActive" class="mx-5">
        <v-tab-item value="詳細情報">
          <div v-for="m in metadata" :key="m.text" class="mb-2 d-flex">
            <span class="d-inline-block vulnerability-label">{{ m.text }}</span>
            <div>
              <span class="mr-3">:</span>
              {{ vulnerability[m.value] }}
            </div>
          </div>
          <div class="mb-2">
            <div class="mb-1">リクエスト</div>
            <pre
              class="pa-3 blue-grey darken-4 white--text w-100 overflow-x-auto"
              >{{ vulnerability.request }}</pre
            >
          </div>
          <div>
            <div class="mb-1">レスポンス</div>
            <pre
              class="pa-3 blue-grey darken-4 white--text w-100 overflow-x-auto"
              >{{ vulnerability.response }}</pre
            >
          </div>
        </v-tab-item>
        <v-tab-item value="コメント">
          <MessagesView
            :load-message-ids="fetchMessageIds"
            :load-messages="fetchMessages"
            :message-ids="messageIds"
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
  pathDastVulnerabilityMessages,
  pathProjects,
  pathVulnerabilities,
  Url
} from "@/assets/js/urls.js";
let user = null;
export default {
  layout: "auth",
  components: {
    MessagesView
  },
  data() {
    return {
      messageIds: [],
      messages: [],
      metadata: [
        {
          text: "パス",
          value: "path"
        },
        {
          text: "名前",
          value: "name"
        },
        {
          text: "メソッド",
          value: "method"
        }
      ],
      tabActive: "",
      tabs: ["詳細情報", "コメント"],
      users: [],
      vulnerability: null
    };
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
      this.setupSocket();
    });
    this.fetchVulnerabilityAndMessageCountAndUsers();
  },
  methods: {
    fetchMessageIds(completion) {
      const vulnerabilityId = this.$route.params.id;
      const url = new Url(pathDastVulnerabilityMessages);
      const data = {
        vulnerabilityId
      };
      ajax.get(url.getIds, data).then(response => {
        this.messageIds = response.data.filter(
          id => !this.messages.find(message => id === message.id)
        );
        completion();
      });
    },
    fetchMessages(ids, completion) {
      const url = new Url(pathDastVulnerabilityMessages);
      const data = {
        ids
      };
      ajax.get(url.ids, data).then(response => {
        const messages = response.data.filter(
          message => !this.messages.find(m => m.id === message.id)
        );
        this.messages = messages.concat(this.messages);
        completion();
      });
    },
    fetchVulnerabilityAndMessageCountAndUsers() {
      const url = new Url(pathVulnerabilities);
      const data = {
        id: this.$route.params.id
      };
      ajax
        .get(url.base, data)
        .then(response => {
          this.vulnerability = response.data[0];
          const url = new Url(pathDastVulnerabilityMessages);
          const data = {
            vulnerabilityId: this.vulnerability.id
          };
          return ajax.get(url.count, data);
        })
        .then(response => {
          this.messageCount = response.data;
          const url = new Url(pathProjects);
          const data = {
            id: this.vulnerability.scan.projectId
          };
          return ajax.get(url.base, data);
        })
        .then(response => {
          this.users = response.data[0].users;
        });
    },
    setupSocket() {
      const url = new Url(pathDastVulnerabilityMessages);
      const socket = new WebSocket(url.socket(user.id));
      const that = this;
      socket.onmessage = e => {
        if (!that.vulnerability) {
          return;
        }
        const message = JSON.parse(e.data);
        if (message.vulnerabilityId !== that.vulnerability.id) {
          return;
        }
        if (message.userId === user.id) {
          return;
        }
        this.messages.push(message);
      };
    },
    storeMessage(message, parent, completion) {
      const url = new Url(pathDastVulnerabilityMessages);
      const data = {
        text: message,
        vulnerabilityId: this.vulnerability.id,
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
.vulnerability-label {
  width: 100px;
}
</style>
