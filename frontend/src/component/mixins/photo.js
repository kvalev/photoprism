import {Photo, TypeLive, TypeRaw, TypeVideo} from "model/photo";
import Thumb from "model/thumb";

export default {
  data() {
    return {
      results: [],
      selection: [],
      filter: {},
      loading: false,
      scrollDisabled: false,
      batchSize: Photo.batchSize(),
      offset: 0,
      page: 0,
      viewer: {
        results: [],
        loading: false,
      },
    };
  },
  methods: {
    loadMore() {
      if (this.scrollDisabled) return;

      this.scrollDisabled = true;
      this.loading = true;

      const count = this.dirty ? (this.page + 2) * this.batchSize : this.batchSize;
      const offset = this.dirty ? 0 : this.offset;

      const params = {
        count: count,
        offset: offset,
        merged: true,
      };

      Object.assign(params, this.filter);

      Photo.search(params).then(response => {
        this.results = Photo.mergeResponse(this.results, response);

        this.complete = (response.count < count);
        this.scrollDisabled = this.complete;

        if (this.complete) {
          this.offset = offset;

          if (this.results.length > 1) {
            this.$notify.info(this.$gettextInterpolate(this.$gettext("Showing all %{n} results"), {n: this.results.length}));
          }
        } else if (this.results.length >= Photo.limit()) {
          this.offset = offset;
          this.complete = true;
          this.scrollDisabled = true;
          this.$notify.warn(this.$gettext("Can't load more, limit reached"));
        } else {
          this.offset = offset + count;
          this.page++;

          this.$nextTick(() => {
            if (this.$root.$el.clientHeight <= window.document.documentElement.clientHeight + 300) {
              this.$emit("scrollRefresh");
            }
          });
        }
      }).catch(() => {
        this.scrollDisabled = false;
      }).finally(() => {
        this.dirty = false;
        this.loading = false;
      });
    },
    openLocation(index) {
      const photo = this.results[index];

      if (photo.CellID && photo.CellID !== "zz") {
        this.$router.push({name: "place", params: {q: photo.CellID}});
      } else if (photo.PlaceID && photo.PlaceID !== "zz") {
        this.$router.push({name: "place", params: {q: photo.PlaceID}});
      } else if (photo.Country && photo.Country !== "zz") {
        this.$router.push({name: "place", params: {q: "country:" + photo.Country}});
      } else {
        this.$notify.warn("unknown location");
      }
    },
    editPhoto(index) {
      let selection = this.results.map((p) => {
        return p.getId();
      });

      // Open Edit Dialog
      Event.publish("dialog.edit", {selection: selection, album: null, index: index});
    },
    openPhoto(index, showMerged) {
      if (this.loading || this.viewer.loading || !this.results[index]) {
        return false;
      }

      const selected = this.results[index];

      // Don't open as stack when user is selecting pictures, or a RAW has only one JPEG.
      if (this.selection.length > 0 || selected.Type === TypeRaw && selected.jpegFiles().length < 2) {
        showMerged = false;
      }

      if (showMerged && selected.Type === TypeLive || selected.Type === TypeVideo) {
        if (selected.isPlayable()) {
          this.$viewer.play({video: selected, album: this.album});
        } else {
          this.$viewer.show(Thumb.fromPhotos(this.results), index);
        }
      } else if (showMerged) {
        this.$viewer.show(Thumb.fromFiles([selected]), 0);
      } else {
        this.viewerResults().then((results) => {
          const thumbsIndex = results.findIndex(result => result.UID === selected.UID);

          if (thumbsIndex < 0) {
            this.$viewer.show(Thumb.fromPhotos(this.results), index);
          } else {
            this.$viewer.show(Thumb.fromPhotos(results), thumbsIndex);
          }
        });
      }

      return true;
    },
    viewerResults() {
      if (this.loading || this.viewer.loading) {
        return Promise.resolve(this.results);
      }

      if (this.viewer.results.length > (this.results.length + this.batchSize)) {
        return Promise.resolve(this.viewer.results);
      }

      this.viewer.loading = true;

      const count = this.batchSize * (this.page + 6);
      const offset = 0;

      const params = {
        count: count,
        offset: offset,
        merged: true,
      };

      Object.assign(params, this.filter);

      return Photo.search(params).then((resp) => {
        // Success.
        this.viewer.loading = false;
        this.viewer.results = resp.models;
        return Promise.resolve(this.viewer.results);
      }, () => {
        // Error.
        this.viewer.loading = false;
        return Promise.resolve(this.results);
      });
    },
  }
};