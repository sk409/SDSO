<template>
  <v-container v-if="vulnerability">
    <div class="title">脆弱性の詳細</div>
    <v-divider class="mb-7"></v-divider>
    <div v-for="m in metadata" :key="m.text" class="mb-5">
      <span class="d-inline-block vulnerability-label">{{ m.text }}:</span>
      <span>{{ vulnerability[m.value] }}</span>
    </div>
    <v-card>
      <v-tabs v-model="tabActive">
        <v-tab v-for="tab in tabs" :key="tab" :href="'#' + tab">{{
          tab
        }}</v-tab>
      </v-tabs>
      <v-tabs-items v-model="tabActive">
        <v-tab-item
          v-for="tabItem in tabItems"
          :key="tabItem.text"
          :value="tabItem.text"
        >
          <pre class="pa-3">{{ vulnerability[tabItem.value] }}</pre>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { pathVulnerabilities, Url } from "@/assets/js/urls.js";
export default {
  layout: "auth",
  data() {
    return {
      metadata: [
        {
          text: "パス",
          value: "path"
        },
        {
          text: "名前",
          value: "name"
        },
        {
          text: "メソッド",
          value: "method"
        }
      ],
      tabActive: "",
      tabs: ["リクエスト", "レスポンス"],
      tabItems: [
        {
          text: "リクエスト",
          value: "request"
        },
        {
          text: "レスポンス",
          value: "response"
        }
      ],
      vulnerability: null
    };
  },
  created() {
    const url = new Url(pathVulnerabilities);
    const data = {
      id: this.$route.params.id
    };
    ajax.get(url.base, data).then(response => {
      this.vulnerability = response.data[0];
    });
  }
};
</script>

<style>
.vulnerability-label {
  width: 100px;
}
</style>
