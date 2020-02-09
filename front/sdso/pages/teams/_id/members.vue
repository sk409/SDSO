<template>
  <v-container v-if="user" fluid>
    <v-subheader>チーム内のメンバ一覧</v-subheader>
    <v-divider></v-divider>
    <v-row>
      <v-col cols="3" class="text-center">
        <v-btn
          color="primary"
          :disabled="!manager"
          outlined
          @click="dialogs.invitation = true"
          >ユーザを招待</v-btn
        >
        <div v-if="!manager" class="red--text caption">権限がありません</div>
      </v-col>
    </v-row>
    <v-card>
      <v-simple-table>
        <thead>
          <tr>
            <th></th>
            <th>名前</th>
            <th>権限</th>
            <th>参加日</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td :style="{ width: profileImageSize() }">
              <v-avatar :size="profileImageSize()">
                <v-img :src="$serverUrl(user.profileImagePath)"></v-img>
              </v-avatar>
            </td>
            <td>{{ user.name }}</td>
            <td>{{ user.role.role | role }}</td>
            <td>{{ user.joinedAt | dateDefault }}</td>
          </tr>
        </tbody>
      </v-simple-table>
    </v-card>
    <v-dialog v-model="dialogs.invitation">
      <TeamUserInvitation
        :team="team"
        @cancel="dialogs.invitation = false"
      ></TeamUserInvitation>
    </v-dialog>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import roles from "@/assets/js/roles.js";
import TeamUserInvitation from "@/components/TeamUserInvitation.vue";
import { pathTeams, pathTeamUsers, pathUsers, Url } from "@/assets/js/urls.js";
export default {
  layout: "team",
  components: {
    TeamUserInvitation
  },
  data() {
    return {
      dialogs: {
        invitation: false
      },
      team: null,
      user: null,
      users: []
    };
  },
  created() {
    this.$fetchUser().then(response => {
      this.user = response.data;
    });
    this.fetchTeam();
    this.fetchUsers();
  },
  computed: {
    manager() {
      if (!this.user) {
        return false;
      }
      const user = this.users.find(user => user.id === this.user.id);
      return user ? user.role.role === roles.team.manager : false;
    },
    selectedRole() {
      const role = this.dialogs.invitation.roles.find(role => role.checked);
      return role;
    }
  },
  methods: {
    fetchTeam() {
      const teamId = this.$route.params.id;
      const url = new Url(pathTeams);
      const data = {
        id: teamId
      };
      ajax.get(url.base, data).then(response => {
        this.team = response.data[0];
      });
    },
    fetchUsers() {
      const teamId = this.$route.params.id;
      const url = new Url(pathTeamUsers);
      const data = {
        teamId
      };
      ajax.get(url.base, data).then(response => {
        const teamUsers = response.data;
        const userIds = teamUsers.map(teamUser => teamUser.userId);
        const url = new Url(pathUsers);
        const data = {
          ids: userIds
        };
        ajax.get(url.ids, data).then(response => {
          for (const index in response.data) {
            const user = response.data[index];
            user.role = teamUsers[index].role;
            user.joinedAt = teamUsers[index].createdAt;
            this.users.push(user);
          }
        });
      });
    },
    profileImageSize() {
      return "64px";
    }
  }
};
</script>
