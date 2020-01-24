<template>
  <v-container v-if="test" class="mt-4">
    <div class="title">テストの詳細</div>
    <v-divider class="mb-5"></v-divider>
    <div class="mb-5">
      <v-chip :color="test.color" text-color="white" class="title">
        {{ test.status }}
      </v-chip>
    </div>
    <v-card class="pb-6">
      <v-tabs v-model="tabActive" class="mb-3">
        <v-tab v-for="tab in tabs" :key="tab" :href="`#${tab}`">
          {{ tab }}
        </v-tab>
      </v-tabs>
      <v-tabs-items v-model="tabActive" class="mx-3">
        <v-tab-item value="コマンド">
          <v-expansion-panels flat multiple :accordian="false">
            <v-expansion-panel
              v-for="result in test.results"
              :key="result.id"
              :style="{ 'border-left': `5px solid ${result.color}` }"
              class="my-2"
              style="border: 1px solid lightgray;"
            >
              <v-expansion-panel-header class="body-1">
                {{ result.command }}
              </v-expansion-panel-header>
              <v-expansion-panel-content>
                <pre class="black white--text pa-2 console-output">{{
                  result.output
                }}</pre>
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { pathTests, Url } from "@/assets/js/urls.js";
export default {
  layout: "auth",
  data() {
    return {
      tabActive: "",
      tabs: ["コマンド", "コメント"],
      test: null
    };
  },
  created() {
    const url = new Url(pathTests);
    const data = {
      id: this.$route.params.id
    };
    ajax.get(url.base, data).then(response => {
      this.test = response.data[0];
    });
  }
};
</script>
