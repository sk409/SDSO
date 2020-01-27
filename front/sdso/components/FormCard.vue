<template>
  <div>
    <v-container fluid class="mt-5">
      <v-row justify="center">
        <v-col cols="12" md="8">
          <v-card class="pa-4">
            <v-card-title>{{title}}</v-card-title>
            <v-card-text>
              <slot name="form"></slot>
            </v-card-text>
            <v-card-actions>
              <slot name="buttons"></slot>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: ""
    }
  }
};
</script>
