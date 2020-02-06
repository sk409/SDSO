<template>
  <div class="h-100">
    <MainView>
      <template v-slot:sidemenu>
        <div class="d-flex meetings-toolbar">
          <v-btn
            v-if="$store.state.projects.project"
            icon
            class="ml-auto"
            @click="dialog = true"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </div>
        <v-list>
          <v-list-item
            v-for="meeting in meetings"
            :key="meeting.id"
            @click="selectMeeting(meeting)"
          >
            <v-list-item-title>
              {{ meeting.name }}
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
                <v-list-item
                  v-for="user in selectedMeeting.users"
                  :key="user.id"
                >
                  <v-list-item-avatar>
                    <v-img :src="$serverUrl(user.profileImagePath)"></v-img>
                  </v-list-item-avatar>
                  <v-list-item-title>
                    {{ user.name }}
                  </v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </div>
        </div>
        <MessagesView
          :messages="messages"
          @send="sendMessage"
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
      dialog: false,
      meetings: [],
      messages: [],
      selectedMeeting: null,
      users: []
    };
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
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
          this.meetings = meetings.filter(
            meeting => meeting.projectId === project.id
          );
        });
    },
    fetchMessages() {
      const url = new Url(pathMeetingMessages);
      const data = {
        meetingId: this.selectedMeeting.id
      };
      ajax.get(url.base, data).then(response => {
        this.messages = response.data;
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
    sendMessage(message, parent) {
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
        this.messages.push(response.data);
      });
    },
    selectMeeting(meeting) {
      this.selectedMeeting = meeting;
      this.fetchMessages();
    }
  }
};
</script>

<style>
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
