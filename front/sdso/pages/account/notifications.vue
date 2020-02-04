<template>
  <v-container fluid>
    <v-subheader>チームへの招待</v-subheader>
    <v-divider></v-divider>
    <v-row
      v-for="teamUserInvitationRequest in teamUserInvitationRequests"
      :key="teamUserInvitationRequest.id"
      justify="center"
    >
      <v-col cols="8">
        <v-card>
          <v-card-text>
            <div class="mb-5">
              <div class="border-left subtitle-1 mb-2">
                <span class="ml-3">チーム名</span>
              </div>
              <div class="body-1">
                {{ teamUserInvitationRequest.team.name }}
              </div>
            </div>
            <div class="mb-5">
              <div class="border-left subtitle-1 mb-2">
                <span class="ml-3">役割</span>
              </div>
              <div class="body-1">
                {{ teamUserInvitationRequest.role.role | role }}
              </div>
            </div>
            <div class="mb-5">
              <div class="border-left subtitle-1 mb-2">
                <span class="ml-3">メッセージ</span>
              </div>
              <pre class="body-1">{{ teamUserInvitationRequest.message }}</pre>
            </div>
            <div class="mb-5">
              <div class="border-left subtitle-1 mb-2">
                <span class="ml-3">招待者</span>
              </div>
              <div>
                <v-avatar size="64">
                  <v-img
                    :src="
                      $serverUrl(
                        teamUserInvitationRequest.inviterUser.profileImagePath
                      )
                    "
                  ></v-img>
                </v-avatar>
                <span class="ml-3">{{
                  teamUserInvitationRequest.inviterUser.name
                }}</span>
              </div>
            </div>
            <div class="mb-5">
              <div class="border-left subtitle-1 mb-2">
                <span class="ml-3">参加するプロジェクト</span>
              </div>
              <v-list>
                <v-list-item
                  v-for="project in teamUserInvitationRequest.projects"
                  :key="project.id"
                >
                  <span class="mr-3">・</span>{{ project.name }}
                </v-list-item>
              </v-list>
            </div>
            <div class="mb-5">
              <div class="border-left subtitle-1 mb-2">
                <span class="ml-3">招待日</span>
              </div>
              <div class="body-1">
                {{ teamUserInvitationRequest.createdAt | dateDefault }}
              </div>
            </div>
          </v-card-text>
          <v-card-actions>
            <v-btn color="error" large text class="ml-auto">拒否</v-btn>
            <v-btn
              color="primary"
              large
              text
              class="mr-auto"
              @click="
                joinTeam(
                  teamUserInvitationRequest.role.role,
                  teamUserInvitationRequest.teamId,
                  teamUserInvitationRequest.inviteeUserId
                )
              "
              >参加</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
    <v-snackbar v-model="snackbar" :timeout="2000" top>
      <span>{{ notification }}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import {
  pathProjectUsers,
  pathTeams,
  pathTeamUserInvitationRequests,
  pathTeamUsers,
  pathUsers,
  Url
} from "@/assets/js/urls.js";
let user = null;
export default {
  layout: "account",
  data() {
    return {
      notification: "",
      snackbar: false,
      teamUserInvitationRequests: []
    };
  },
  created() {
    this.$fetchUser().then(response => {
      user = response.data;
      this.fetchTeamUserInvitationReuqests();
    });
  },
  methods: {
    fetchTeamUserInvitationReuqests() {
      const url = new Url(pathTeamUserInvitationRequests);
      const data = {
        inviteeUserId: user.id
      };
      ajax.get(url.base, data).then(response => {
        this.teamUserInvitationRequests = response.data;
      });
    },
    joinTeam(role, teamId, userId) {
      const url = new Url(pathTeamUsers);
      const data = {
        role,
        teamId,
        userId
      };
      ajax.post(url.base, data).then(response => {
        const teamUserInvitationRequest = this.teamUserInvitationRequests.find(
          teamUserInvitationRequest =>
            teamUserInvitationRequest.teamId === teamId
        );
        this.teamUserInvitationRequests = this.teamUserInvitationRequests.filter(
          t => t.id !== teamUserInvitationRequest.id
        );
        this.snackbar = true;
        this.notification = `チーム${teamUserInvitationRequest.team.name}に参加しました`;
        for (const project of teamUserInvitationRequest.projects) {
          const url = new Url(pathProjectUsers);
          const data = {
            projectId: project.id,
            userId: user.id,
            role: teamUserInvitationRequest.role.role
          };
          ajax.post(url.base, data);
        }
        const url = new Url(pathTeamUserInvitationRequests);
        const data = {
          id: teamUserInvitationRequest.id
        };
        ajax.delete(url.base, data);
      });
    }
  }
};
</script>

<style>
.border-left {
  border-left: 5px solid #3f51b5;
}
</style>
