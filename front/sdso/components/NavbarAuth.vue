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
              hide-details
              :items="projectnames"
              no-data-text="プロジェクトがありません"
              :value="activeProjectName"
              label="プロジェクトを選択してください"
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
export default {
  props: {
    activeProjectName: {
      type: String,
      default: ""
    },
    height: {
      type: Number,
      default: 64
    },
    projects: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      drawer: false,
      menuItems: [
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
          title: "脆弱性",
          icon: "mdi-security",
          route: ""
        },
        {
          title: "アカウント情報",
          icon: "mdi-account-badge-horizontal",
          route: ""
        }
      ]
    };
  },
  computed: {
    projectnames() {
      return this.projects.map(project => project.name);
    }
  },
  methods: {
    selectProjectname(projectname) {
      this.$emit("select-projectname", projectname);
    }
  }
};
</script>
