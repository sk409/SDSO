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
            <v-row>
              <v-col cols="2" class="subtitle-1">チーム名</v-col>
              <v-col>{{teamUserInvitationRequest.team.name}}</v-col>
            </v-row>
            <v-row>
              <v-col cols="2" class="subtitle-1">招待日</v-col>
              <v-col>{{teamUserInvitationRequest.createdAt | dateDefault}}</v-col>
            </v-row>
            <v-row>
              <v-col cols="2" class="subtitle-1">招待者</v-col>
            </v-row>
            <v-row align="center">
              <v-col cols="3">
                <v-img :src="$serverUrl(teamUserInvitationRequest.user.profileImagePath)"></v-img>
              </v-col>
              <v-col>{{teamUserInvitationRequest.user.name}}</v-col>
            </v-row>
          </v-card-text>
          <v-card-actions>
            <v-btn color="error" text class="ml-auto">拒否</v-btn>
            <v-btn
              color="primary"
              text
              class="mr-auto"
              @click="joinTeam(teamUserInvitationRequest.role, teamUserInvitationRequest.teamId, teamUserInvitationRequest.userId)"
            >参加</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import {
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
        userId: user.id
      };
      ajax.get(url.base, data).then(response => {
        const teamUserInvitationRequests = response.data;
        const teamIds = teamUserInvitationRequests.map(
          teamUserInvitationRequest => teamUserInvitationRequest.teamId
        );
        const url = new Url(pathTeams);
        const data = {
          ids: teamIds
        };
        ajax.get(url.ids, data).then(response => {
          const teams = response.data;
          const userIds = teamUserInvitationRequests.map(
            teamUserInvitationRequest => teamUserInvitationRequest.userId
          );
          const url = new Url(pathUsers);
          const data = {
            ids: userIds
          };
          ajax.get(url.ids, data).then(response => {
            const users = response.data;
            for (const index in teamUserInvitationRequests) {
              teamUserInvitationRequests[index].team = teams.find(
                team => team.id === teamUserInvitationRequests[index].teamId
              );
              teamUserInvitationRequests[index].user = users.find(
                user => user.id === teamUserInvitationRequests[index].userId
              );
            }
            this.teamUserInvitationRequests = teamUserInvitationRequests;
          });
        });
      });
    },
    joinTeam(role, teaId, userId) {
      const url = new Url(pathTeamUsers);
      const data = {
        role,
        teamId,
        userId
      };
    }
  }
};
</script>