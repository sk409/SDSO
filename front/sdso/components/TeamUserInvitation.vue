<template>
  <v-card>
    <v-stepper v-model="step">
      <v-stepper-header flat>
        <v-stepper-step :complete="step > 1" step="1">権限選択</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step :complete="step > 2" step="2">招待</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step step="3">完了</v-stepper-step>
      </v-stepper-header>
      <v-stepper-items>
        <v-stepper-content step="1">
          <div class="pl-2">
            <v-subheader>招待するユーザの権限を指定してください</v-subheader>
            <div v-for="role in roles" :key="role.text">
              <div class="d-flex align-center mb-3">
                <v-checkbox
                  :input-value="role.checked"
                  :label="role.text"
                  :error="role.error"
                  :error-messages="role.errorMessages"
                  @change="checkRole(role.value)"
                ></v-checkbox>
              </div>
              <div>
                <div class="caption">{{role.caption}}</div>
                <v-list>
                  <v-list-item v-for="ability in role.ablities" :key="ability">{{ability}}</v-list-item>
                </v-list>
              </div>
            </div>
          </div>
          <v-divider class="my-2"></v-divider>
          <div class="d-flex justify-end">
            <v-btn text @click="$emit('cancel')">キャセル</v-btn>
            <v-btn text color="primary" @click="selectRole">次へ</v-btn>
          </div>
        </v-stepper-content>
        <v-stepper-content step="2">
          <v-row>
            <v-col cols="7">
              <v-text-field
                v-model="username"
                :error="usernameError"
                :error-messages="usernameErrorMessages"
                label="招待するメンバの名前を入力してエンターを押してください"
                @keydown.enter="inviteUser"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row class="mt-0">
            <v-col cols="3">
              <v-list>
                <v-list-item v-for="invitedUser in invitedUsers" :key="invitedUser.id">
                  <v-list-item-avatar>
                    <v-img :src="$serverUrl(invitedUser.profileImagePath)"></v-img>
                  </v-list-item-avatar>
                  <v-list-item-content>
                    <v-list-item-title>{{invitedUser.name}}</v-list-item-title>
                  </v-list-item-content>
                  <v-list-item-action>
                    <v-btn icon @click="removeUser(invitedUser.name)">
                      <v-icon>mdi-delete</v-icon>
                    </v-btn>
                  </v-list-item-action>
                </v-list-item>
              </v-list>
            </v-col>
          </v-row>
          <v-row class="mt-0">
            <v-col cols="6">
              <v-textarea v-model="message" :rules="messageRules" label="メッセージ"></v-textarea>
            </v-col>
          </v-row>
          <div class="d-flex justify-end">
            <v-btn text @click="$emit('cancel')">キャセル</v-btn>
            <v-btn text color="primary" @click="step=1">戻る</v-btn>
            <v-btn text color="primary" @click="selectUsers">次へ</v-btn>
          </div>
        </v-stepper-content>
        <v-stepper-content step="3">
          <v-subheader>以下のメンバを{{selectedRole.text}}として招待しました</v-subheader>
          <v-list>
            <v-list-item v-for="invitedUser in invitedUsers" :key="invitedUser.id">
              <v-list-item-avatar>
                <v-img :src="$serverUrl(invitedUser.profileImagePath)"></v-img>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title>{{invitedUser.name}}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <div class="d-flex justify-end">
            <v-btn color="primary" text @click="close">閉じる</v-btn>
          </div>
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </v-card>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import roles from "@/assets/js/roles.js";
import {
  pathTeamUserInvitationRequests,
  pathUsers,
  Url
} from "@/assets/js/urls.js";
let messageMaxLength = 512;
let user = null;
export default {
  props: {
    team: {
      required: true,
      validator: v => typeof v === "object" || v === null
    }
  },
  data() {
    return {
      invitedUsers: [],
      message: "",
      messageRules: [
        v =>
          v.length <= messageMaxLength ||
          `${messageMaxLength}文字以内で入力してください`
      ],
      projects: [],
      role: roles.team.manager,
      roles: [
        {
          ablities: [
            "・チーム内のメンバの管理権限",
            "・チーム内のプロジェクトの管理権限",
            "・チームにメンバーを招待する権限"
          ],
          caption: "一般ユーザの権限に加えて、以下の権限が与えられます。",
          checked: true,
          error: false,
          errorMessages: [],
          text: "管理者",
          value: roles.team.manager
        },
        {
          ablities: [
            "・チーム内の許可されたリソースへのアクセス権限",
            "・プロジェクトの作成権限"
          ],
          caption: "以下の権限が与えられます。",
          checked: false,
          error: false,
          errorMessages: [],
          text: "一般ユーザ",
          value: roles.team.member
        }
      ],
      step: 1,
      username: "",
      usernameError: false,
      usernameErrorMessages: [],
      visible: false
    };
  },
  computed: {
    selectedRole() {
      return this.roles.find(role => role.checked);
    }
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
    });
  },
  methods: {
    checkRole(role) {
      this.role = role;
      for (const r of this.roles) {
        if (r.value === role) {
          r.checked = !r.checked;
        } else {
          r.checked = false;
        }
      }
    },
    close() {
      this.step = 1;
      this.invitedUsers = [];
      this.message = "";
      this.$emit("cancel");
    },
    inviteUser(e) {
      this.usernameError = false;
      this.usernameErrorMessages = [];
      const username = e.target.value;
      const url = new Url(pathUsers);
      const data = {
        name: username
      };
      ajax.get(url.base, data).then(response => {
        if (response.data.length === 0) {
          return;
        }
        this.username = "";
        this.invitedUsers.push(response.data[0]);
      });
    },
    removeUser(username) {
      this.invitedUsers = this.invitedUsers.filter(invitedUser => {
        return false;
        return invitedUser.name !== username;
      });
    },
    selectRole() {
      const selectedRole = this.roles.find(role => role.checked);
      if (!selectedRole) {
        this.roles.forEach(role => {
          role.error = true;
          role.errorMessages = ["どちらかの権限を選択してください"];
        });
        return;
      }
      this.step = 2;
    },
    selectUsers() {
      let error = false;
      if (this.invitedUsers.length === 0) {
        error = true;
        this.usernameError = true;
        this.usernameErrorMessages = "ユーザを招待してください";
      }
      if (messageMaxLength < this.message.length) {
        error = true;
      }
      if (error) {
        return;
      }
      const url = new Url(pathTeamUserInvitationRequests);
      for (const invitedUser of this.invitedUsers) {
        const data = {
          message: this.message,
          teamId: this.team.id,
          userId: user.id
        };

        ajax.post(url.base, data);
      }
      this.step = 3;
    }
  }
};
</script>