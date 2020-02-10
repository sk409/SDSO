<template>
  <v-simple-table>
    <thead>
      <tr>
        <th>パス</th>
        <th>種類</th>
        <th>メソッド</th>
        <th>SHA1</th>
        <th>状態</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="vulnerability in vulnerabilities"
        :key="vulnerability.id"
        @click="$router.push($routes.vulnerabilities.show(vulnerability.id))"
      >
        <td>{{ vulnerability.path }}</td>
        <td>{{ vulnerability.name }}</td>
        <td>{{ vulnerability.method }}</td>
        <td>{{ vulnerability.scan.commitSha1.substr(0, 7) }}</td>
        <td>
          <v-chip color="red" small text-color="white">未修正</v-chip>
        </td>
      </tr>
    </tbody>
  </v-simple-table>
</template>

<script>
export default {
  props: {
    vulnerabilities: {
      type: Array,
      default: []
    }
  }
};
</script>
