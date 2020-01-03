<template>
  <div>
    <div id="header-container">
      <n-link to="/" id="app-title">SDSO</n-link>
      <div id="header-right">
        <!-- <div id="logout-button" @click="logout">ログアウト</div> -->
        <!-- <div class="separator">|</div> -->
        <!-- <div></div> -->
        <div class="hamburger" @click="showMenu"></div>
      </div>
    </div>
    <transition @after-leave="hideMenuCompletion">
      <div v-if="menu.isVisible" class="menu w-100">
        <div class="container-fluid h-100">
          <div class="row h-100">
            <div class="h-100 col-12 col-lg-6 offset-lg-6 menu-col">
              <div class="mt-2 d-flex">
                <i
                  class="close-menu-button ml-auto mr-2 el-icon-close bg-white p-1 rounded-circle"
                  @click="hideMenu"
                ></i>
              </div>
              <div class="text-center">
                <img class="profile-image" :src="profileImagePath" />
              </div>
              <div class="text-white text-center h3">{{ user.Name }}</div>
              <div class="border-bottom border-white"></div>
              <div class="w-50 mx-auto mt-4 text-center text-white">
                <div class="p-2 menu-button" @click="transition($routes.dashboardProjects)">ダッシュボード</div>
                <div class="mt-3 p-2 menu-button" @click="logout">サインアウト</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  name: "AppHeader",
  data() {
    return {
      menu: {
        isVisible: false
      },
      user: null,
      hideMenuCompletionHandler: null
    };
  },
  computed: {
    profileImagePath() {
      return this.user.ProfileImagePath
        ? this.user.ProfileImagePath
        : "@/assets/images/no-image.png";
    }
  },
  created() {
    this.$ajax.get(this.$urls.user, {}, { withCredentials: true }, response => {
      this.user = response.data;
    });
  },
  methods: {
    logout() {
      this.$ajax.post(
        this.$urls.logout,
        {},
        { withCredentials: true },
        response => {
          if (response.status !== 200) {
            this.$notify.error({
              message: "ログアウトに失敗しました",
              duration: 3000
            });
            return;
          }
          this.$notify.success({
            message: "ログアウトしました",
            duration: 3000
          });
          this.$router.push(this.$routes.login);
        }
      );
    },
    showMenu() {
      this.menu.isVisible = true;
    },
    hideMenu() {
      this.menu.isVisible = false;
    },
    hideMenuCompletion() {
      if (this.hideMenuCompletionHandler) {
        this.hideMenuCompletionHandler();
        this.hideMenuCompletionHandler = null;
      }
    },
    transition(route) {
      this.hideMenuCompletionHandler = () => {
        this.$router.push(this.$routes.dashboardProjects);
      };
      this.hideMenu();
    }
  }
};
</script>

<style scoped>
#header-container {
  display: flex;
  align-items: center;
  width: 100%;
  height: 100%;
  color: rgb(212, 212, 212);
  padding: 0 1rem;
}

#app-title {
  display: flex;
  align-items: center;
  color: rgb(212, 212, 212);
  font-size: 1.35rem;
}

#header-right {
  display: flex;
  align-items: center;
  font-size: 1rem;
  margin-left: auto;
}

#logout-button:hover {
  cursor: pointer;
  text-decoration: underline;
}

.separator {
  font-weight: bold;
  margin: 0 0.8rem;
}

.hamburger {
  background: rgb(112, 112, 112);
  transform: rotate(90deg);
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  position: relative;
  cursor: pointer;
}

.hamburger::after {
  content: "|||";
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.hamburger:hover {
  background: rgb(132, 132, 132);
}

.menu {
  position: absolute;
  z-index: 5;
  height: 100%;
  right: 0;
  top: 0;
  transition: all 0.5s;
}

.menu-col {
  background: #2196f3;
}

.close-menu-button {
  width: 2rem;
  height: 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.profile-image {
  width: 200px;
  height: 200px;
  object-fit: cover;
  border-radius: 50%;
}

.menu-button {
  cursor: pointer;
  border: 1px solid white;
}

.menu-button:hover {
  border: 1px solid #e91e63;
}

.v-enter,
.v-leave-to {
  transform: translateX(100%);
}
.v-enter-to,
.v-leave {
  transform: translateX(0);
}
</style>
