<template>
  <div>
    <v-subheader>動的テスト結果一覧</v-subheader>
    <v-divider class="mb-3"></v-divider>
    <GitToolbar class="mb-3" @change-revision="fetchScans"></GitToolbar>
    <v-row justify="center">
      <v-col cols="11">
        <v-card class="mb-4">
          <v-simple-table>
            <thead>
              <tr>
                <th>パス</th>
                <th>種類</th>
                <th>メソッド</th>
                <th>状態</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="vulnerability in vulnerabilities"
                :key="vulnerability.id"
                @click="
                  $router.push($routes.vulnerabilities.show(vulnerability.id))
                "
              >
                <td>{{ vulnerability.path }}</td>
                <td>{{ vulnerability.name }}</td>
                <td>{{ vulnerability.method }}</td>
                <td>
                  <v-chip color="red" small text-color="white">未修正</v-chip>
                </td>
              </tr>
            </tbody>
          </v-simple-table>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import GitToolbar from "@/components/GitToolbar.vue";
import { pathScans, Url } from "@/assets/js/urls.js";
export default {
  layout: "dashboard",
  components: {
    GitToolbar
  },
  data() {
    return {
      scans: [],
      user: null
    };
  },
  computed: {
    vulnerabilities() {
      return this.scans.map(scan => scan.vulnerabilities).flat();
    }
  },
  created() {
    this.$nuxt.$emit("setSidemenuType", "security");
    this.$fetchUser().then(response => {
      this.user = response.data;
      this.fetchScans();
    });
  },
  methods: {
    fetchScans() {
      const project = this.$store.state.projects.project;
      if (!project) {
        return;
      }
      const revision = this.$store.state.git.revision;
      if (!revision) {
        return;
      }
      const url = new Url(pathScans);
      const data = {
        projectname: project.name,
        revision,
        username: this.user.name
      };
      ajax.get(url.base, data).then(response => {
        this.scans = response.data;
      });
    }
  }
};
</script>
