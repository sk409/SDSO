<template>
  <MainView>
    <template v-slot:sidemenu>
      <v-list class="pa-2">
        <v-list-item
          v-for="sidemenuItem in sidemenuItems"
          :key="sidemenuItem.title"
          :to="sidemenuItem.route"
          router
        >
          <v-list-item-action>
            <v-icon>{{ sidemenuItem.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ sidemenuItem.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </template>
    <template v-slot:content>
      <v-subheader>動的テスト結果一覧</v-subheader>
      <v-divider class="mb-1"></v-divider>
      <GitToolbar class="mb-3" @change-revision="fetchScans"></GitToolbar>
      <v-row v-if="vulnerabilities.length" justify="center">
        <v-col cols="11">
          <v-card class="mb-4">
            <v-simple-table>
              <thead>
                <tr>
                  <th>パス</th>
                  <th>種類</th>
                  <th>メソッド</th>
                  <th>状態</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="vulnerability in vulnerabilities"
                  :key="vulnerability.id"
                  @click="
                    $router.push($routes.vulnerabilities.show(vulnerability.id))
                  "
                >
                  <td>{{ vulnerability.path }}</td>
                  <td>{{ vulnerability.name }}</td>
                  <td>{{ vulnerability.method }}</td>
                  <td>
                    <v-chip color="red" small text-color="white">未修正</v-chip>
                  </td>
                </tr>
              </tbody>
            </v-simple-table>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </MainView>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import GitToolbar from "@/components/GitToolbar.vue";
import MainView from "@/components/MainView.vue";
import { pathScans, Url } from "@/assets/js/urls.js";
export default {
  layout: "dashboard",
  components: {
    GitToolbar,
    MainView
  },
  data() {
    return {
      scans: [],
      sidemenuItems: [
        {
          title: "動的テスト",
          icon: "mdi-shield-plus-outline",
          route: this.$routes.dashboard.dast
        },
        {
          title: "パッケージ",
          icon: "mdi-package-variant",
          route: ""
        }
      ],
      user: null
    };
  },
  computed: {
    vulnerabilities() {
      return this.scans.map(scan => scan.vulnerabilities).flat();
    }
  },
  created() {
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchScans();
    });
  },
  methods: {
    fetchScans() {
      const team = this.$store.state.teams.team;
      if (!team) {
        return;
      }
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const revision = this.$store.state.git.revision;
      if (!revision) {
        return;
      }
      const url = new Url(pathScans);
      const data = {
        projectname: project.name,
        revision,
        teamname: team.name,
        username: this.user.name
      };
      ajax.get(url.base, data).then(response => {
        this.scans = response.data;
      });
    }
  }
};
</script>
