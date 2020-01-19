<template>
  <v-app light>
    <NavbarAuth></NavbarAuth>
    <v-content class="white black--text">
      <v-container fluid class="h-100">
        <v-layout class="h-100">
          <v-flex xs0 md2 class="d-none d-md-block sidemenu pr-2">
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
                  <v-list-item-title>
                    {{ sidemenuItem.title }}
                  </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-flex>
          <v-flex xs12 md10>
            <nuxt />
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import NavbarAuth from "@/components/NavbarAuth.vue";
export default {
  middleware: "auth",
  components: {
    NavbarAuth
  },
  data() {
    return {
      sidemenuItems: []
    };
  },
  created() {
    this.$nuxt.$on("setSidemenuItems", this.setSidemenuItems);
  },
  methods: {
    setSidemenuItems(sidemenuItems) {
      this.sidemenuItems = sidemenuItems;
    }
  }
};
</script>

<style>
.sidemenu {
  border-right: 2px solid lightgrey;
}
</style>
