<template>
  <div v-infinite-scroll="loadMore" class="p-tab p-tab-discover-today" :infinite-scroll-disabled="scrollDisabled"
       :infinite-scroll-distance="1200" :infinite-scroll-listen-for-event="'scrollRefresh'">

    <v-container v-if="loading" fluid class="pa-4">
      <v-progress-linear color="secondary-dark" :indeterminate="true"></v-progress-linear>
    </v-container>
    <v-container v-else-if="!results.length" text-xs-center fluid class="pa-4">
      <p class="subheading pb-3"><translate>No photos on this day.</translate></p>
    </v-container>
    <v-container v-else fluid class="pa-0">
      <p-scroll-top></p-scroll-top>

      <div v-for="year in years" :key="year">
        <h3 class="pl-3 pt-1">{{ year }}</h3>

        <!-- No clipboard, no toolbar only a single, card-based view (no edit-photo, no open-location, no selection) -->
        <p-photo-cards context="photos"
                      :photos="photosByYear(year)"
                      :disable-selection="true"
                      :filter="filter"
                      :open-photo="openPhotoByYear(year)"></p-photo-cards>
      </div>
    </v-container>
  </div>
</template>

<script>
import PhotoMixin from "component/mixins/photo.js";
import {TypeLive, TypeRaw, TypeVideo} from "model/photo";
import Thumb from "model/thumb";

export default {
  name: 'PTabDiscoverToday',
  mixins: [PhotoMixin],
  data() {
    const today = new Date();

    return {
      filter: {
        month: today.getMonth() + 1,
        day: today.getDate(),
        order: 'newest',
      }
    };
  },
  computed: {
    years: function() {
      const yrs = new Set()

      for (let i = 0; i < this.results.length; i++) {
        let item = this.results[i];
        yrs.add(item.Year);
      }

      return yrs;
    }
  },
  methods: {
    photosByYear(year) {
      const photos = [];

      for (let item of this.results) {
        if (item.Year == year) {
          photos.push(item);
        }
      }

      return photos;
    },
    photosOffsetByYear(selectedYear) {
      let offset = 0;

      for (let year of this.years) {
        if (year > selectedYear) {
          offset += this.photosByYear(year).length;
        }
      }

      return offset;
    },
    openPhotoByYear(year) {
      return (index, showMerged) => {
        const offset = this.photosOffsetByYear(year);
        index += offset;

        this.openPhoto(index, showMerged);
      };
    },
  }
};
</script>
