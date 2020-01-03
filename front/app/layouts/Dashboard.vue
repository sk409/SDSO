<template>
  <div class="h-100">
    <AppHeader id="app-header"></AppHeader>
    <div id="app-body" class="d-flex">
      <div class="sidebar d-none d-lg-block">
        <n-link
          :to="$routes.dashboardProjects"
          :style="sidebarMenuStyle($routes.dashboardProjects)"
          class="sidebar-menu"
          >プロジェクト一覧</n-link
        >
      </div>
      <div class="w-100">
        <nuxt />
      </div>
      <div class="bottom-menu d-block d-lg-none">
        <el-select v-model="bottomMenuValue">
          <el-option label="プロジェクト一覧" value="project"> </el-option>
        </el-select>
      </div>
    </div>
  </div>
</template>

<script>
import AppHeader from "@/components/AppHeader.vue";
export default {
  name: "Dashboard",
  //middleware: "auth",
  components: {
    AppHeader
  },
  data() {
    return {
      bottomMenuValue: "project"
    };
  },
  methods: {
    sidebarMenuStyle(route) {
      //console.log(this.$route);
      if (this.$route.path === route) {
        return {
          color: "rgb(200, 178, 64)"
        };
      }
      return {};
    }
  }
};
</script>

<style scoped>
.sidebar {
  width: 15%;
  height: 100%;
  background: rgb(41, 43, 48);
  padding-top: 0.8rem;
}

.sidebar-menu {
  color: rgb(212, 212, 212);
  text-align: center;
  width: 100%;
  display: block;
}

.bottom-menu {
  position: absolute;
  bottom: 36px;
  left: 50%;
  transform: translateX(-55%);
  box-shadow: 0px 0px 10px 1px rgba(0, 0, 0, 0.3);
}
</style>
