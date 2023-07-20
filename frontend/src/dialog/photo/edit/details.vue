<template>
  <div class="p-tab p-tab-photo-details">
    <v-container fluid>
      <v-form ref="form" lazy-validation
              dense class="p-form-photo-details-meta" accept-charset="UTF-8"
              @submit.prevent="save">
        <v-layout row wrap align-top fill-height>
          <v-flex
            class="p-photo pa-2"
            xs12 sm4 md2
          >
            <v-card tile
                    class="ma-0 elevation-0"
                    :title="model.Title">
              <v-img v-touch="{left, right}"
                     :src="model.thumbnailUrl('tile_500')"
                     aspect-ratio="1"
                     class="card darken-1 elevation-0 clickable"
                     @click.exact="openPhoto()"
              >
              </v-img>

            </v-card>
          </v-flex>
          <v-flex xs12 sm8 md10 fill-height>
            <v-layout row wrap>
              <v-flex xs12 lg6 class="pa-2">
                <v-text-field
                  v-model="model.Title"
                  :append-icon="model.TitleSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :rules="[textRule]"
                  hide-details box flat
                  :label="$pgettext('Photo', 'Title')"
                  placeholder=""
                  color="secondary-dark"
                  browser-autocomplete="off"
                  class="input-title"
                ></v-text-field>
              </v-flex>

              <v-flex xs4 md1 pa-2>
                <v-autocomplete
                  v-model="model.Day"
                  :append-icon="model.TakenSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :error="invalidDate"
                  :label="$gettext('Day')"
                  browser-autocomplete="off"
                  hide-details box flat hide-no-data
                  color="secondary-dark"
                  :items="options.Days()"
                  class="input-day"
                  @change="updateTime">
                </v-autocomplete>
              </v-flex>
              <v-flex xs4 md1 pa-2>
                <v-autocomplete
                  v-model="model.Month"
                  :append-icon="model.TakenSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :error="invalidDate"
                  :label="$gettext('Month')"
                  browser-autocomplete="off"
                  hide-details box flat hide-no-data
                  color="secondary-dark"
                  :items="options.MonthsShort()"
                  class="input-month"
                  @change="updateTime">
                </v-autocomplete>
              </v-flex>
              <v-flex xs4 md2 pa-2>
                <v-autocomplete
                  v-model="model.Year"
                  :append-icon="model.TakenSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :error="invalidDate"
                  :label="$gettext('Year')"
                  browser-autocomplete="off"
                  hide-details box flat hide-no-data
                  color="secondary-dark"
                  :items="options.Years()"
                  class="input-year"
                  @change="updateTime">
                </v-autocomplete>
              </v-flex>

              <v-flex xs6 md2 class="pa-2">
                <v-text-field
                  v-model="time"
                  :append-icon="model.TakenSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :label="model.timeIsUTC() ? $gettext('Time UTC') : $gettext('Local Time')"
                  browser-autocomplete="off"
                  autocorrect="off"
                  autocapitalize="none"
                  hide-details box flat
                  return-masked-value mask="##:##:##"
                  color="secondary-dark"
                  class="input-local-time"
                ></v-text-field>
              </v-flex>

              <v-flex xs6 sm6 md6 lg3 class="pa-2">
                <v-autocomplete
                  v-model="model.TimeZone"
                  :disabled="disabled"
                  :label="$gettext('Time Zone')"
                  browser-autocomplete="off"
                  hide-details box flat hide-no-data
                  color="secondary-dark"
                  item-value="ID"
                  item-text="Name"
                  :items="options.TimeZones()"
                  class="input-timezone"
                  @change="updateTime">
                </v-autocomplete>
              </v-flex>

              <v-flex xs12 sm8 md4 lg3 class="pa-2">
                <v-autocomplete
                  v-model="model.Country"
                  :append-icon="model.PlaceSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :readonly="!!(model.Lat || model.Lng)"
                  :label="$gettext('Country')" hide-details box flat
                  hide-no-data
                  browser-autocomplete="off"
                  color="secondary-dark"
                  item-value="Code"
                  item-text="Name"
                  :items="countries"
                  class="input-country">
                </v-autocomplete>
              </v-flex>

              <v-flex xs4 md2 lg2 class="pa-2">
                <v-text-field
                  v-model="model.Altitude"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  autocorrect="off"
                  autocapitalize="none"
                  :label="$gettext('Altitude (m)')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-altitude"
                ></v-text-field>
              </v-flex>

              <v-flex xs4 sm6 md3 lg2 class="pa-2">
                <v-text-field
                  v-model="model.Lat"
                  :append-icon="model.PlaceSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  autocorrect="off"
                  autocapitalize="none"
                  :label="$gettext('Latitude')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-latitude"
                ></v-text-field>
              </v-flex>

              <v-flex xs4 sm6 md3 lg2 class="pa-2">
                <v-text-field
                  v-model="model.Lng"
                  :append-icon="model.PlaceSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  autocorrect="off"
                  autocapitalize="none"
                  :label="$gettext('Longitude')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-longitude"
                ></v-text-field>
              </v-flex>

              <v-flex xs12 sm12 md12>
                <div id="map" style="height: 500px;" />
              </v-flex>

              <v-flex xs12 md6 pa-2 class="p-camera-select">
                <v-select
                  v-model="model.CameraID"
                  :append-icon="model.CameraSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :label="$gettext('Camera')"
                  :menu-props="{'maxHeight':346}"
                  browser-autocomplete="off"
                  hide-details box flat
                  color="secondary-dark"
                  item-value="ID"
                  item-text="Name"
                  :items="cameraOptions"
                  class="input-camera">
                </v-select>
              </v-flex>

              <v-flex xs6 md3 class="pa-2">
                <v-text-field
                  v-model="model.Iso"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  autocorrect="off"
                  autocapitalize="none"
                  label="ISO"
                  placeholder=""
                  color="secondary-dark"
                  class="input-iso"
                ></v-text-field>
              </v-flex>

              <v-flex xs6 md3 class="pa-2">
                <v-text-field
                  v-model="model.Exposure"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  autocorrect="off"
                  autocapitalize="none"
                  :label="$gettext('Exposure')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-exposure"
                ></v-text-field>
              </v-flex>

              <v-flex xs12 md6 pa-2 class="p-lens-select">
                <v-select
                  v-model="model.LensID"
                  :append-icon="model.CameraSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :label="$gettext('Lens')"
                  :menu-props="{'maxHeight':346}"
                  browser-autocomplete="off"
                  hide-details box flat
                  color="secondary-dark"
                  item-value="ID"
                  item-text="Name"
                  :items="lensOptions"
                  class="input-lens">
                </v-select>
              </v-flex>

              <v-flex xs6 md3 class="pa-2">
                <v-text-field
                  v-model="model.FNumber" f
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  autocorrect="off"
                  autocapitalize="none"
                  :label="$gettext('F Number')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-fnumber"
                ></v-text-field>
              </v-flex>

              <v-flex xs6 md3 class="pa-2">
                <v-text-field
                  v-model="model.FocalLength"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  :label="$gettext('Focal Length')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-focal-length"
                ></v-text-field>
              </v-flex>

              <v-flex xs12 md6 class="pa-2">
                <v-text-field
                  v-model="model.Details.Artist"
                  :append-icon="model.Details.ArtistSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :rules="[textRule]"
                  hide-details box flat
                  browser-autocomplete="off"
                  :label="$gettext('Artist')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-artist"
                ></v-text-field>
              </v-flex>

              <v-flex xs6 md3 class="pa-2">
                <v-text-field
                  v-model="model.Details.Copyright"
                  :append-icon="model.Details.CopyrightSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :rules="[textRule]"
                  hide-details box flat
                  browser-autocomplete="off"
                  :label="$gettext('Copyright')"
                  placeholder=""
                  color="secondary-dark"
                  class="input-copyright"
                ></v-text-field>
              </v-flex>

              <v-flex xs6 md3 class="pa-2">
                <v-textarea
                  v-model="model.Details.License"
                  :append-icon="model.Details.LicenseSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :rules="[textRule]"
                  hide-details box flat
                  browser-autocomplete="off"
                  auto-grow
                  :label="$gettext('License')"
                  placeholder=""
                  :rows="1"
                  color="secondary-dark"
                  class="input-license"
                ></v-textarea>
              </v-flex>

              <v-flex xs12 class="pa-2">
                <v-textarea
                  v-model="model.Details.Subject"
                  :append-icon="model.Details.SubjectSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  :rules="[textRule]"
                  hide-details box flat
                  browser-autocomplete="off"
                  auto-grow
                  :label="$gettext('Subject')"
                  placeholder=""
                  :rows="1"
                  color="secondary-dark"
                  class="input-subject"
                ></v-textarea>
              </v-flex>

              <v-flex xs12 class="pa-2">
                <v-textarea
                  v-model="model.Description"
                  :append-icon="model.DescriptionSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  auto-grow
                  :label="$gettext('Description')"
                  placeholder=""
                  :rows="1"
                  color="secondary-dark"
                  class="input-description"
                ></v-textarea>
              </v-flex>

              <v-flex xs12 md8 class="pa-2">
                <v-textarea
                  v-model="model.Details.Keywords"
                  :append-icon="model.Details.KeywordsSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  auto-grow
                  :label="$gettext('Keywords')"
                  placeholder=""
                  :rows="1"
                  color="secondary-dark"
                  class="input-keywords"
                ></v-textarea>
              </v-flex>

              <v-flex xs12 md4 class="pa-2">
                <v-textarea
                  v-model="model.Details.Notes"
                  :append-icon="model.Details.NotesSrc === 'manual' ? 'check' : ''"
                  :disabled="disabled"
                  hide-details box flat
                  browser-autocomplete="off"
                  auto-grow
                  :label="$gettext('Notes')"
                  placeholder=""
                  :rows="1"
                  color="secondary-dark"
                  class="input-notes"
                ></v-textarea>
              </v-flex>

              <v-flex v-if="!disabled" xs12 :text-xs-right="!rtl" :text-xs-left="rtl" class="pt-3">
                <v-btn depressed color="secondary-light" class="compact action-close"
                       @click.stop="close">
                  <translate>Close</translate>
                </v-btn>
                <v-btn color="primary-button" depressed dark class="compact action-apply action-approve"
                       @click.stop="save(false)">
                  <span v-if="$config.feature('review') && model.Quality < 3"><translate>Approve</translate></span>
                  <span v-else><translate>Apply</translate></span>
                </v-btn>
                <v-btn color="primary-button" depressed dark class="compact action-done hidden-xs-only"
                       @click.stop="save(true)">
                  <translate>Done</translate>
                </v-btn>
              </v-flex>
            </v-layout>
          </v-flex>
        </v-layout>

        <div class="mt-1 clear"></div>
      </v-form>
    </v-container>
  </div>
</template>

<script>
import maplibregl from "maplibre-gl";
import MapboxDraw from "@mapbox/mapbox-gl-draw";
import MaplibreGeocoder from '@maplibre/maplibre-gl-geocoder';

import "@mapbox/mapbox-gl-draw/dist/mapbox-gl-draw.css";
import '@maplibre/maplibre-gl-geocoder/dist/maplibre-gl-geocoder.css';

import { LocateNearbyControl } from "common/map-controls";
import nominatim from "common/nominatim";
import Thumb from "model/thumb";
import Photo from "model/photo";
import countries from "options/countries.json";
import * as options from "options/options";

export default {
  name: 'PTabPhotoDetails',
  props: {
    model: {
      type: Object,
      default: () => new Photo(false),
    },
    uid: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      disabled: !this.$config.feature("edit"),
      config: this.$config.values,
      all: {
        colors: [{label: this.$gettext("Unknown"), name: ""}],
      },
      readonly: this.$config.get("readonly"),
      options: options,
      countries: countries,
      showDatePicker: false,
      showTimePicker: false,
      invalidDate: false,
      time: "",
      textRule: v => v.length <= this.$config.get('clip') || this.$gettext("Text too long"),
      rtl: this.$rtl,
      map: null,
    };
  },
  computed: {
    cameraOptions() {
      return this.config.cameras;
    },
    lensOptions() {
      return this.config.lenses;
    },
  },
  watch: {
    model() {
      this.updateTime();
    },
    uid() {
      this.updateTime();
    },
  },
  created() {
    this.updateTime();
  },
  mounted() {
    this.renderMap();
  },
  methods: {
    updateTime() {
      if (!this.model.hasId()) {
        return;
      }

      const taken = this.model.getDateTime();

      this.time = taken.toFormat("HH:mm:ss");
    },
    updateModel() {
      if (!this.model.hasId()) {
        return;
      }

      let localDate = this.model.localDate(this.time);

      this.invalidDate = !localDate.isValid;

      if (this.invalidDate) {
        return;
      }

      if (this.model.Day === 0) {
        this.model.Day = parseInt(localDate.toFormat("d"));
      }

      if (this.model.Month === 0) {
        this.model.Month = parseInt(localDate.toFormat("L"));
      }

      if (this.model.Year === 0) {
        this.model.Year = parseInt(localDate.toFormat("y"));
      }

      const isoTime = localDate.toISO({
        suppressMilliseconds: true,
        includeOffset: false,
      }) + "Z";

      this.model.TakenAtLocal = isoTime;

      if (this.model.currentTimeZoneUTC()) {
        this.model.TakenAt = isoTime;
      }
    },
    left() {
      this.$emit('next');
    },
    right() {
      this.$emit('prev');
    },
    openPhoto() {
      this.$viewer.show(Thumb.fromFiles([this.model]), 0);
    },
    save(close) {
      if (this.invalidDate) {
        this.$notify.error(this.$gettext("Invalid date"));
        return;
      }

      this.updateModel();

      this.model.update().then(() => {
        if (close) {
          this.$emit('close');
        }

        this.updateTime();
      });
    },
    close() {
      this.$emit('close');
    },
    renderMap() {
      const s = this.config.settings.maps;

      let mapKey = "";

      if (this.$config.has("mapKey")) {
        mapKey = this.$config.get("mapKey");
      }

      let mapOptions = {
        container: "map",
        cooperativeGestures: true,
        style: "https://api.maptiler.com/maps/" + s.style + "/style.json?key=" + mapKey,
        attributionControl: true,
        customAttribution: "<a href='https://nominatim.org/' target='_blank'>Geocoding provided by OpenStreetMap Nominatim</a>",
        zoom: 0,
      };

      const localStorageGeocoder = (query) => {
        const indicator = '<span class="material-icons md-12">cached</span>';

        const matches = [];
        const geolocations = JSON.parse(window.localStorage.getItem("geolocations")) || [];

        for (const location of geolocations) {
          if (location.place_name.toLowerCase().includes(query.toLowerCase())) {
            location.place_name = `${indicator} ${location.place_name}`;
            matches.push(location);
          }
        }

        return matches;
      };
      const geocoder = new MaplibreGeocoder(nominatim, {
        maplibregl: maplibregl,
        placeholder: this.$gettext("Search"),
        marker: true,
        localGeocoder: localStorageGeocoder,
        // the explicit search is broken in several ways:
        // - the enter key event is retargeted to the clear button, so no search is performed and instead the text is cleared
        // - even if the above is fixed, the keyup.enter event is bubbled-up to the form and a save action is performed
        showResultsWhileTyping: true,
        debounceSearch: 1500,
        limit: 10,
        // strip the cached indicator when showing the place in the search bar
        getItemValue: (item) => item.place_name.split("</span>").pop().trim(),
        flyTo: {
          duration: this.config.settings.maps.animate,
          essential: false,
          animate: this.config.settings.maps.animate > 0,
        }
      });

      geocoder.on('result', (e) => {
        // if the user selects a cached entry (meaning a previously selected location), instead of only jumping
        // to that location, we should also set it as the photo's geolocation
        if (e.result.properties.cached) {
          draw.add(e.result.geometry);
          this.map.fire('draw.create', { features: [e.result]});
        }
      });

      const draw = new MapboxDraw({
        displayControlsDefault: false,
        controls: {
            point: true,
            trash: true
        },
      });

      this.map = new maplibregl.Map(mapOptions);
      this.map.setLanguage(this.config.settings.ui.language.split("-")[0]);

      const controlPos = this.$rtl ? 'top-left' : 'top-right';

      this.map.addControl(geocoder, controlPos);
      this.map.addControl(new LocateNearbyControl({
        model: this.model,
        settings: this.config.settings.maps,
        range: { hours: 1 },
      }), controlPos);
      this.map.addControl(new maplibregl.NavigationControl({showCompass: true}), controlPos);
      this.map.addControl(new maplibregl.FullscreenControl({container: document.querySelector('body')}), controlPos);
      this.map.addControl(new maplibregl.GeolocateControl({
        positionOptions: {
          enableHighAccuracy: true
        },
        trackUserLocation: true
      }), controlPos);

      this.map.addControl(draw, controlPos);

      this.map.on('draw.create', async (e) => {
        // remove all previously selected points
        const featureCollection = draw.getAll();
        const previous = featureCollection.features.slice(0, -1);
        previous.forEach(feature => draw.delete(feature.id));

        // update the model with the selected location
        const point = e.features[0].geometry.coordinates;
        this.model.Lat = point[1];
        this.model.Lng = point[0];

        // in case of a cached location we can skip the reverse lookup and the storing
        // in local storage, as this has already been done
        if (e.properties && e.properties.cached) {
          return;
        }

        // reverse geocode the selected location and store it in localStorage
        const lookupConfig = {query: [this.model.Lat, this.model.Lng]};
        const lookupResult = await nominatim.reverseGeocode(lookupConfig);

        if (lookupResult && lookupResult.features && lookupResult.features.length) {
          // set a marker that the location has been cached
          const geolocation = lookupResult.features[0];
          geolocation.properties.cached = true;

          const geolocations = JSON.parse(window.localStorage.getItem("geolocations")) || [];

          if (geolocations.some((g) => g.id === geolocation.id)) {
            return;
          }

          geolocations.push(geolocation);
          window.localStorage.setItem("geolocations", JSON.stringify(geolocations));
        }
      });
      this.map.on('draw.delete', e => {
        // for some reason the backend does not reset the country code whenever the location is unset,
        // so let's do it in the frontend
        // this will force the country estimation (based on the path) to kick in
        this.model.Country = "zz";
        this.model.Lat = 0;
        this.model.Lng = 0;
      });
      this.map.on('draw.update', e => {
        if (e.action === "move") {
          // update the model with the updated location
          const point = e.features[0].geometry.coordinates;
          this.model.Lat = point[1];
          this.model.Lng = point[0];
        }
      });
      this.map.on("radar.selected", e => {
        this.model.Lat = e.location.lat;
        this.model.Lng = e.location.lng;
      });

      this.map.on('load', () => {
        if (this.model.hasLocation()) {
          draw.add({ type: 'Point', coordinates: [this.model.Lng, this.model.Lat] });
          this.map.flyTo({
            center: [this.model.Lng, this.model.Lat],
            zoom: 16,
            maxDuration: 0.1,
          });
        }
      });
    },
  },
};
</script>
