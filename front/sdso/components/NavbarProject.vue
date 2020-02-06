<template>
  <div>
    <v-navigation-drawer v-model="drawer" app>
      <v-list>
        <v-list-item
          v-for="navItem in navItems"
          :key="navItem.title"
          :to="navItem.route"
          router
        >
          <v-list-item-action>
            <v-icon>{{ navItem.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ navItem.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar color="primary" app dark fixed>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-container fluid>
        <v-row align="center">
          <v-col cols="3">
            <v-select
              v-model="teamname"
              hide-details
              :items="teamnames"
              no-data-text="チームがありません"
              placeholder="チームを選択してください"
              @input="selectTeamname"
            ></v-select>
          </v-col>
          <v-col cols="3">
            <v-select
              v-model="projectname"
              hide-details
              :items="projectnames"
              no-data-text="プロジェクトがありません"
              placeholder="プロジェクトを選択してください"
              @input="selectProjectname"
            ></v-select>
          </v-col>
          <v-col cols="1" offset="5" class="d-flex align-center">
            <v-btn
              icon
              small
              class="ml-3"
              @click="$router.push($routes.account.teams)"
            >
              <v-icon>mdi-account</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-app-bar>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import { mapMutations } from "vuex";
import {
  pathProjects,
  pathProjectUsers,
  pathTeams,
  pathTeamUsers,
  Url
} from "@/assets/js/urls.js";

let unsubscribe = null;
let user = null;
export default {
  data() {
    return {
      drawer: false,
      navItems: [
        {
          title: "コード管理",
          icon: "mdi-source-branch",
          route: this.$routes.dashboard.commits.base
        },
        {
          title: "テスト",
          icon: "mdi-test-tube",
          route: this.$routes.dashboard.tests
        },
        {
          title: "セキュリティ",
          icon: "mdi-security",
          route: this.$routes.dashboard.dast
        },
        {
          title: "メンバ",
          icon: "mdi-account-multiple-outline",
          route: this.$routes.dashboard.members
        },
        {
          title: "ミーティング",
          icon: "mdi-message-outline",
          route: this.$routes.dashboard.meetings
        }
      ],
      projectname: "",
      projects: [],
      teamname: "",
      teams: []
    };
  },
  computed: {
    projectnames() {
      return this.projects.map(project => project.name);
    },
    teamnames() {
      return this.teams.map(team => team.name);
    }
  },
  created() {
    const team = this.$store.state.teams.team;
    if (team) {
      this.teamname = team.name;
    }
    const project = this.$store.state.projects.project;
    if (project) {
      this.projectname = project.name;
    }
    this.$fetchUser().then(response => {
      user = response.data;
      this.fetchTeams();
      this.fetchProjects();
    });
    unsubscribe = this.$store.subscribe(mutation => {
      if (mutation.type !== mutations.teams.setTeam) {
        return;
      }
      this.setProject(null);
      this.fetchProjects();
    });
  },
  destroyed() {
    unsubscribe();
  },
  methods: {
    ...mapMutations({
      setProject: mutations.projects.setProject,
      setTeam: mutations.teams.setTeam
    }),
    fetchTeams() {
      const url = new Url(pathTeamUsers);
      const data = {
        userId: user.id
      };
      ajax
        .get(url.base, data)
        .then(response => {
          const teamIds = response.data.map(teamUser => teamUser.teamId);
          return teamIds;
        })
        .then(teamIds => {
          if (teamIds.length === 0) {
            return;
          }
          const url = new Url(pathTeams);
          const data = {
            ids: teamIds
          };
          ajax.get(url.ids, data).then(response => {
            this.teams = response.data;
            this.teams = this.teams.sort((a, b) => (a.name < b.name ? -1 : 1));
          });
        });
    },
    fetchProjects() {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
      const url = new Url(pathProjects);
      const data = {
        teamId: team.id
      };
      ajax.get(url.base, data).then(response => {
        this.projects = response.data;
      });
    },
    selectProjectname(projectname) {
      const project = this.projects.find(
        project => project.name === projectname
      );
      if (!project) {
        return;
      }
      this.setProject(project);
    },
    selectTeamname(teamname) {
      const team = this.teams.find(team => team.name === teamname);
      if (!team) {
        return;
      }
      this.setTeam(team);
    }
  }
};
</script>
