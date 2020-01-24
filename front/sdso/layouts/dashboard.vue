<template>
  <v-app light>
    <NavbarAuth
      ref="navbar"
      :active-project-name="activeProjectName"
      :projects="projects"
      :height="navbarHeight"
      @select-projectname="selectProjectname"
    ></NavbarAuth>
    <v-content class="white black--text h-100">
      <div ref="content" class="d-flex h-100">
        <div class="h-100 sidemenu">
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
        </div>
        <div
          :style="mainStyle"
          class="w-100 h-100 overflow-x-hidden overflow-y-auto"
        >
          <nuxt />
        </div>
      </div>
      <!-- <v-container fluid class="h-100">
        <v-row class="h-100">
          <v-col md="3" class="d-none d-md-block sidemenu pr-2">
            <v-list>
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
                  <v-list-item-title>{{
                    sidemenuItem.title
                  }}</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-col>
          <v-col cols="12" md="9">
            <nuxt />
          </v-col>
        </v-row>
      </v-container>-->
    </v-content>
  </v-app>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mutations from "@/assets/js/mutations.js";
import NavbarAuth from "@/components/NavbarAuth.vue";
import { mapMutations } from "vuex";
import { pathProjects, Url } from "@/assets/js/urls.js";
export default {
  middleware: "auth",
  components: {
    NavbarAuth
  },
  data() {
    return {
      allSidemenuItems: {
        git: [
          {
            title: "コミット",
            icon: "mdi-source-commit",
            route: this.$routes.dashboard.commits
          },
          {
            title: "ファイル",
            icon: "mdi-file-outline",
            route: this.$routes.dashboard.files()
          }
        ],
        security: [
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
        tests: [
          {
            title: "テスト結果",
            icon: "mdi-test-tube",
            route: this.$routes.dashboard.tests
          }
        ]
      },
      mainStyle: {},
      navbarHeight: 64,
      projects: [],
      sidemenuType: "git",
      user: null
    };
  },
  computed: {
    activeProjectName() {
      const project = this.$store.state.projects.project;
      return project ? project.name : "";
    },
    sidemenuItems() {
      return this.allSidemenuItems[this.sidemenuType];
    }
  },
  created() {
    this.$nuxt.$on("setSidemenuType", this.setSidemenuType);
    this.$fetchUser()
      .then(response => {
        this.user = response.data;
        const url = new Url(pathProjects);
        const data = {
          userId: this.user.id
        };
        return ajax.get(url.base, data);
      })
      .then(response => {
        this.projects = response.data;
      });
  },
  mounted() {
    let maxHeight = this.$refs.content.clientHeight;
    if (this.$refs.navbar.$el.clientHeight === 0) {
      maxHeight += this.navbarHeight;
    }
    this.mainStyle = {
      "max-height": maxHeight + "px"
    };
    console.log(this.mainStyle);
  },
  methods: {
    ...mapMutations({
      setProject: mutations.projects.setProject
    }),
    selectProjectname(projectname) {
      const url = new Url(pathProjects);
      const data = {
        name: projectname,
        userId: this.user.id
      };
      ajax.get(url.base, data).then(response => {
        const project = response.data[0];
        this.setProject(project);
      });
    },
    setSidemenuType(sidemenuType) {
      this.sidemenuType = sidemenuType;
    }
  }
};
</script>

<style>
.sidemenu {
  border-right: 2px solid lightgrey;
}
</style>
