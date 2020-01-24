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
    <v-app-bar color="primary" :height="height" app dark fixed>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-container fluid>
        <v-row align="center">
          <v-col cols="4">
            <v-select
              v-model="projectname"
              hide-details
              :items="projectnames"
              no-data-text="プロジェクトがありません"
              placeholder="プロジェクトを選択してください"
              class="project-select"
              @input="selectProjectname"
            ></v-select>
          </v-col>
          <v-col cols="1" offset="7">
            <v-menu offset-y>
              <template v-slot:activator="{ on }">
                <v-btn icon small v-on="on">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item
                  v-for="menuItem in menuItems"
                  :key="menuItem.title"
                  :to="menuItem.route"
                  router
                >
                  <v-list-item-action>
                    <v-icon>{{ menuItem.icon }}</v-icon>
                  </v-list-item-action>
                  <v-list-item-content>
                    <v-list-item-title>{{ menuItem.title }}</v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </v-menu>
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
import { pathProjects, Url } from "@/assets/js/urls.js";

let user = null;
export default {
  props: {
    height: {
      type: Number,
      default: 64
    }
  },
  data() {
    return {
      drawer: false,
      menuItems: [
        {
          title: "チーム作成",
          icon: "mdi-account-supervisor-outline",
          route: ""
        },
        {
          title: "プロジェクト作成",
          icon: "mdi-apps",
          route: this.$routes.projects.create
        }
      ],
      navItems: [
        {
          title: "コード管理",
          icon: "mdi-source-branch",
          route: this.$routes.dashboard.commits
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
          title: "ミーティング",
          icon: "mdi-message-outline",
          route: ""
        },
        {
          title: "設定",
          icon: "mdi-settings-outline",
          route: ""
        }
      ],
      projectname: "",
      projects: []
    };
  },
  computed: {
    projectnames() {
      return this.projects.map(project => project.name);
    }
  },
  created() {
    const project = this.$store.state.projects.project;
    if (project) {
      this.projectname = project.name;
    }
    this.$fetchUser().then(response => {
      user = response.data;
      this.fetchProjects();
    });
  },
  methods: {
    ...mapMutations({
      setProject: mutations.projects.setProject
    }),
    selectProjectname(projectname) {
      const project = this.projects.find(
        project => project.name === projectname
      );
      if (!project) {
        return;
      }
      this.setProject(project);
    },
    fetchProjects() {
      const url = new Url(pathProjects);
      const data = {
        userId: user.id
      };
      ajax.get(url.base, data).then(response => {
        this.projects = response.data;
      });
    }
  }
};
</script>

<style>
.project-select {
  width: 200px;
}
</style>
