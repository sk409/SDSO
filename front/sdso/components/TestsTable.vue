<template>
  <v-simple-table>
    <thead>
      <tr>
        <th>ステータス</th>
        <th>ブランチ</th>
        <th>SHA1</th>
        <th>ステップ</th>
        <th>実施日</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="test in tests"
        :key="test.id"
        @click="$router.push($routes.tests.show(test.id))"
      >
        <td>
          <v-chip :color="test.color" small text-color="white">
            {{ test.status }}
          </v-chip>
        </td>
        <td>{{ test.branchname }}</td>
        <td>{{ test.commitSha1.substr(0, 5) }}</td>
        <td>{{ test.results.length }}/{{ test.steps }}</td>
        <td>{{ test.createdAt | dateDefault }}</td>
      </tr>
    </tbody>
  </v-simple-table>
</template>

<script>
export default {
  props: {
    tests: {
      type: Array,
      default: []
    }
  }
};
</script>
