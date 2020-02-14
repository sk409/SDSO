<template>
  <div class="h-100">
    <MainView>
      <template v-slot:sidemenu>
        <div class="d-flex meetings-toolbar">
          <v-btn v-if="$store.state.projects.project" icon class="ml-auto" @click="dialog = true">
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </div>
        <v-list>
          <v-list-item
            v-for="meeting in meetings"
            :key="meeting.id"
            @click="selectMeeting(meeting)"
          >
            <v-list-item-title class="d-flex">
              <span>{{ meeting.name }}</span>
              <span class="ml-auto">{{meeting.newMessageCount | numberLimit(99)}}</span>
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </template>
      <template v-slot:content v-if="selectedMeeting">
        <div class="messages-toolbar d-flex align-center px-3">
          <span class="title">{{ selectedMeeting.name }}</span>
          <div class="ml-auto">
            <v-menu offset-y>
              <template v-slot:activator="{ on }">
                <v-btn v-on="on" text>
                  <v-icon>mdi-account-multiple-outline</v-icon>
                  <span>{{ selectedMeeting.users.length }}</span>
                  <v-icon>mdi-menu-down</v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item v-for="user in selectedMeeting.users" :key="user.id">
                  <v-list-item-avatar>
                    <v-img :src="$serverUrl(user.profileImagePath)"></v-img>
                  </v-list-item-avatar>
                  <v-list-item-title>{{ user.name }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </div>
        </div>
        <MessagesView
          :baseline-message="baselineMessage"
          :load-message-ids="fetchMessageIds"
          :load-messages="fetchMessages"
          :message-ids="messageIds"
          :messages="messages"
          :post-message="storeMessage"
          :users="users"
          class="messages"
        ></MessagesView>
      </template>
    </MainView>
    <v-dialog v-model="dialog" class="w-75">
      <v-card class="pa-5">
        <MeetingForm
          :project="$store.state.projects.project"
          :users="users"
          @created="createdMeeting"
        ></MeetingForm>
      </v-card>
    </v-dialog>
    <NotificationSnackbar></NotificationSnackbar>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import MainView from "@/components/MainView.vue";
import MeetingForm from "@/components/MeetingForm.vue";
import MessagesView from "@/components/MessagesView.vue";
import mutations from "@/assets/js/mutations.js";
import NotificationSnackbar from "@/components/NotificationSnackbar.vue";
import {
  pathMeetings,
  pathMeetingMessages,
  pathMeetingMessageViewers,
  pathMeetingUsers,
  pathProjectUsers,
  pathUsers,
  Url
} from "@/assets/js/urls.js";
let unsubscribe = null;
let user = null;
export default {
  layout: "dashboard",
  components: {
    MainView,
    MeetingForm,
    MessagesView,
    NotificationSnackbar
  },
  data() {
    return {
      baselineMessage: null,
      dialog: false,
      meetings: [],
      messageIds: [],
      messages: [],
      selectedMeeting: null,
      users: []
    };
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
      this.setupScoket();
      this.fetchMeetings();
      this.fetchUsers();
    });
    unsubscribe = this.$store.subscribe(mutation => {
      if (mutation.type === mutations.projects.setProject) {
        this.fetchMeetings();
        this.fetchUsers();
      }
    });
  },
  destroyed() {
    unsubscribe();
  },
  methods: {
    createdMeeting(meeting) {
      this.dialog = false;
      meeting.newMessages = [];
      this.meetings.push(meeting);
    },
    fetchMeetings() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathMeetingUsers);
      const data = {
        userId: user.id
      };
      ajax
        .get(url.base, data)
        .then(response => {
          const meetingUsers = response.data;
          const meetingIds = meetingUsers.map(
            meetingUser => meetingUser.meetingId
          );
          const url = new Url(pathMeetings);
          const data = {
            ids: meetingIds
          };
          return ajax.get(url.ids, data);
        })
        .then(response => {
          const meetings = response.data;
          meetings.forEach(meeting => {
            meeting.newMessages = [];
            meeting.newMessageCount = 0;
          });
          this.meetings = meetings.filter(
            meeting => meeting.projectId === project.id
          );
          this.meetings.forEach(meeting => {
            const url = new Url(pathMeetingMessages);
            const data = {
              meetingId: meeting.id,
              viewerId: user.id
            };
            ajax.get(url.new, data).then(response => {
              meeting.newMessages = response.data;
              meeting.newMessageCount = response.data.length;
            });
          });
        });
    },
    fetchMessageIds(completion) {
      const url = new Url(pathMeetingMessages);
      const data = {
        meetingId: this.selectedMeeting.id
      };
      ajax.get(url.getIds, data).then(response => {
        this.messageIds = response.data;
        completion();
      });
    },
    fetchMessages(ids, completion) {
      const url = new Url(pathMeetingMessages);
      const data = {
        ids
      };
      ajax.get(url.ids, data).then(response => {
        const messages = response.data;
        messages.forEach(message => {
          if (this.selectedMeeting.newMessages.length != 0) {
            message.new =
              this.selectedMeeting.newMessages.find(
                newMessage => newMessage.id === message.id
              ) !== undefined;
            message.newMessageChip =
              this.selectedMeeting.newMessages[0].id === message.id;
          }
        });
        this.messages = messages.concat(this.messages);
        completion();
      });
    },
    fetchUsers() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const url = new Url(pathProjectUsers);
      const data = {
        projectId: project.id
      };
      ajax
        .get(url.base, data)
        .then(response => {
          const projectUsers = response.data;
          const userIds = projectUsers.map(projectUser => projectUser.userId);
          const url = new Url(pathUsers);
          const data = {
            ids: userIds
          };
          return ajax.get(url.ids, data);
        })
        .then(response => {
          this.users = response.data;
        });
    },
    selectMeeting(meeting) {
      meeting.newMessages.forEach(newMessage => {
        const url = new Url(pathMeetingMessageViewers);
        const data = {
          meetingMessageId: newMessage.id,
          userId: user.id
        };
        ajax.post(url.base, data);
      });
      if (meeting.newMessages.length != 0) {
        this.baselineMessage = meeting.newMessages[0];
      }
      meeting.newMessageCount = 0;
      this.selectedMeeting = meeting;
    },
    setupScoket() {
      const url = new Url(pathMeetingMessages);
      const socket = new WebSocket(url.socket(user.id));
      socket.onmessage = e => {
        const message = JSON.parse(e.data);
        if (
          !this.selectedMeeting ||
          message.meetingId !== this.selectedMeeting.id
        ) {
          const meeting = this.meetings.find(
            meeting => meeting.id === message.meetingId
          );
          if (meeting) {
            meeting.newMessages.push(message);
            meeting.newMessageCount += 1;
          }
          return;
        }
        if (message.user.id === user.id) {
          return;
        }
        this.messages.push(message);
        const url = new Url(pathMeetingMessageViewers);
        const data = {
          meetingMessageId: message.id,
          userId: user.id
        };
        ajax.post(url.base, data);
      };
    },
    storeMessage(message, parent, completion) {
      const url = new Url(pathMeetingMessages);
      const data = {
        text: message,
        meetingId: this.selectedMeeting.id,
        userId: user.id
      };
      if (parent) {
        data.parentId = parent.id;
      }
      ajax.post(url.base, data).then(response => {
        const message = response.data;
        this.messages.push(message);
        const url = new Url(pathMeetingMessageViewers);
        const data = {
          meetingMessageId: message.id,
          userId: user.id
        };
        ajax.post(url.base, data);
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
.meetings-toolbar {
  height: 20px;
}
.messages {
  height: 90%;
}
.messages-toolbar {
  height: 10%;
  border-bottom: 1px solid lightgray;
}
</style>
