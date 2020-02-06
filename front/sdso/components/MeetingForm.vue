<template>
  <div>
    <v-form>
      <v-text-field v-model="meetingname" label="名前"></v-text-field>
      <v-select
        :items="usernames"
        label="参加させるメンバを選択してください"
        @input="selectUser"
      ></v-select>
      <v-container>
        <v-row>
          <v-col cols="3">
            <v-list>
              <v-list-item
                v-for="inviteeUser in inviteeUsers"
                :key="inviteeUser.id"
              >
                <v-list-item-avatar>
                  <v-img
                    :src="$serverUrl(inviteeUser.profileImagePath)"
                  ></v-img>
                </v-list-item-avatar>
                <v-list-item-content>
                  {{ inviteeUser.name }}
                </v-list-item-content>
                <v-list-item-action>
                  <v-icon>mdi-delete</v-icon>
                </v-list-item-action>
              </v-list-item>
            </v-list>
          </v-col>
        </v-row>
      </v-container>
    </v-form>
    <div class="text-center">
      <v-btn color="accent" :loading="creating" @click="create">作成</v-btn>
    </div>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import { pathMeetings, pathMeetingUsers, Url } from "@/assets/js/urls.js";
import { mapMutations } from "vuex";
let user = null;
export default {
  props: {
    project: {
      validator: v => typeof v === "object" || v === null,
      required: true
    },
    users: {
      validator: v => typeof v === "object" || v === null,
      required: true
    }
  },
  data() {
    return {
      creating: false,
      items: [],
      inviteeUsers: [],
      meetingname: ""
    };
  },
  computed: {
    usernames() {
      return this.items.map(user => user.name);
    }
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
    });
  },
  watch: {
    users: {
      immediate: true,
      handler(newValue) {
        this.items = newValue;
      }
    }
  },
  methods: {
    ...mapMutations({
      setMessage: mutations.notifications.setMessage
    }),
    create() {
      const url = new Url(pathMeetings);
      const data = {
        name: this.meetingname,
        projectId: this.project.id
      };
      this.creating = true;
      ajax.post(url.base, data).then(response => {
        const meeting = response.data;
        this.setMessage(`ミーティング「${this.meetingname}」を作成しました。`);
        for (const inviteeUser of this.inviteeUsers) {
          const url = new Url(pathMeetingUsers);
          const data = {
            meetingId: meeting.id,
            userId: inviteeUser.id
          };
          ajax.post(url.base, data);
        }
        this.meetingname = "";
        this.$emit("created", meeting);
        this.creating = false;
        this.inviteeUsers = [];
      });
    },
    selectUser(username) {
      const user = this.items.find(user => user.name === username);
      this.inviteeUsers.push(user);
      this.items = this.items.filter(item => item.id !== user.id);
    }
  }
};
</script>
