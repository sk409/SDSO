<template>
  <div>
    <AppHeader id="app-header"></AppHeader>
    <div id="app-body">
      <div class="project-header">
        <div class="metadata">
          <n-link :to="$routes.dashboardProjects">
            {{ pathParamUserName }}
          </n-link>
          <span>/</span>
          <n-link
            :to="$routes.projectCode(pathParamUserName, pathParamProjectName)"
            >{{ pathParamProjectName }}</n-link
          >
        </div>
        <div class="tabs mt-3 mt-lg-0">
          <div class="spacer"></div>
          <n-link
            tag="div"
            :to="$routes.projectCode(pathParamUserName, pathParamProjectName)"
            :style="
              tabStyle(
                $routes.projectCode(pathParamUserName, pathParamProjectName),
                false
              )
            "
            class="tab"
            >コード</n-link
          >
          <n-link
            tag="div"
            :event="isAuthor ? 'click' : ''"
            :to="$routes.projectTest(pathParamUserName, pathParamProjectName)"
            :style="
              tabStyle(
                $routes.projectTest(pathParamUserName, pathParamProjectName),
                true
              )
            "
            class="tab"
            >テスト</n-link
          >
          <n-link
            tag="div"
            :event="isAuthor ? 'click' : ''"
            :to="
              $routes.projectVulnerabilities(
                pathParamUserName,
                pathParamProjectName
              )
            "
            :style="
              tabStyle(
                $routes.projectVulnerabilities(
                  pathParamUserName,
                  pathParamProjectName
                ),
                true
              )
            "
            class="tab"
            >脆弱性</n-link
          >
          <n-link
            tag="div"
            :event="isAuthor ? 'click' : ''"
            :to="
              $routes.projectSettings(pathParamUserName, pathParamProjectName)
            "
            :style="
              tabStyle(
                $routes.projectSettings(
                  pathParamUserName,
                  pathParamProjectName
                ),
                true
              )
            "
            class="tab"
            >設定</n-link
          >
          <div class="spacer"></div>
        </div>
      </div>
      <div class="p-3">
        <nuxt />
      </div>
    </div>
  </div>
</template>

<script>
import AppHeader from "@/components/AppHeader.vue";
export default {
  middleware: "auth",
  components: {
    AppHeader
  },
  data() {
    return {
      user: null
    };
  },
  computed: {
    pathParamUserName() {
      return this.$route.params.userName;
    },
    pathParamProjectName() {
      return this.$route.params.projectName
        ? this.$route.params.projectName
        : this.$route.params.pathMatch;
    },
    isAuthor() {
      return this.user && this.user.Name === this.pathParamUserName;
    }
  },
  created() {
    this.fetchUser();
  },
  methods: {
    tabStyle(route, applyingToColor) {
      const style = {
        color: !applyingToColor || this.isAuthor ? "black" : "gray"
      };
      if (this.$route.path.startsWith(route)) {
        return Object.assign(style, {
          "border-top": "2px solid rgb(212, 105, 43)",
          "border-left": "1px solid rgb(180, 180, 180)",
          "border-right": "1px solid rgb(180, 180, 180)",
          "border-bottom": "none",
          background: "rgb(245, 245, 245)"
        });
      }
      return style;
    },
    fetchUser() {
      this.$ajax.get(
        this.$urls.user,
        {},
        { withCredentials: true },
        response => {
          this.user = response.data;
        }
      );
    }
  }
};
</script>

<style scoped>
.project-header {
  background: rgb(240, 240, 240);
}
.metadata {
  font-size: 1.25rem;
  margin: 0.6rem 0 0 1rem;
}
.tabs {
  display: flex;
  align-items: flex-end;
}
.tab {
  display: inline-block;
  height: 100%;
  padding: 0.5rem 1rem;
  margin: 0;
  border-bottom: 1px solid rgb(180, 180, 180);
}
.tab:hover {
  cursor: pointer;
  text-decoration: underline;
}
.spacer {
  flex: 1;
  border-bottom: 1px solid rgb(180, 180, 180);
}
</style>
