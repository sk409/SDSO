<template>
  <div>
    <div id="header-container">
      <n-link to="/" id="app-title">SDSO</n-link>
      <div id="header-right">
        <div id="logout-button" @click="logout">ログアウト</div>
        <div class="separator">|</div>
        <div></div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "AppHeader",
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
</style>
