<template>
  <div v-infinite-scroll="loadMore" class="p-tab p-tab-discover-today" :infinite-scroll-disabled="scrollDisabled"
       :infinite-scroll-distance="1200" :infinite-scroll-listen-for-event="'scrollRefresh'">

    <v-container v-if="loading" fluid class="pa-4">
      <v-progress-linear color="secondary-dark" :indeterminate="true"></v-progress-linear>
    </v-container>
    <v-container v-else-if="!results.length" text-xs-center fluid class="pa-4">
      <p class="subheading pb-3"><translate>No photos found</translate></p>
    </v-container>
    <v-container v-else fluid class="pa-0">
      <p-scroll-top></p-scroll-top>

      <!-- No clipboard, no toolbar only a single, card-based view (no edit-photo, no open-location, no selection) -->
      <p-photo-cards context="photos"
                    :photos="results"
                    :select-mode="false"
                    :filter="filter"
                    :open-photo="openPhoto"></p-photo-cards>
    </v-container>
  </div>
</template>

<script>
import PhotoMixin from "component/mixins/photo.js";

export default {
  name: 'PTabDiscoverRandom',
  mixins: [PhotoMixin],
  data() {
    return {
      filter: {
        order: 'random',
      },
    };
  },
};
</script>
