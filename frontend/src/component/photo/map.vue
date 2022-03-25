<template>
   <!--  -->
  <v-container fluid fill-height class="p-page p-page-album-photos pa-0">
    <div id="map" style="width: 100%; height: 100%;">
    </div>
  </v-container>
</template>

<script>
import maplibregl from "maplibre-gl";
import Photo from "model/photo";

export default {
  name: 'PPhotoMap',
  props: {
    photos: Array,
    album: Object,
    context: String,
  },
  data() {
    return {
      map: null,
      markers: {},
      markersOnScreen: {},
      url: "",
      attribution: '<a href="https://www.maptiler.com/copyright/" target="_blank">&copy; MapTiler</a> <a href="https://www.openstreetmap.org/copyright" target="_blank">&copy; OpenStreetMap contributors</a>',
      options: {},
      mapFont: [],
      config: this.$config.values,
      settings: this.$config.values.settings.maps,
    };
  },
  mounted() {
    this.$scrollbar.hide();
    this.configureMap().then(() => this.renderMap());
  },
  destroyed() {
    this.$scrollbar.show();
  },
  methods: {
    configureMap() {
      console.log("configure map");
      return this.$config.load().finally(() => {
        const s = this.$config.values.settings.maps;

        let mapKey = "";

        if (this.$config.has("mapKey")) {
          mapKey = this.$config.get("mapKey");
        }

        let mapFont = ['Roboto', 'sans-serif'];

        let mapOptions = {
          container: "map",
          style: "https://api.maptiler.com/maps/" + s.style + "/style.json?key=" + mapKey,
          attributionControl: true,
          customAttribution: this.attribution,
          zoom: 0,
        };

        if (!mapKey || s.style === "offline") {
          mapFont = ["Open Sans Semibold"];
          mapOptions = {
            container: "map",
            style: {
              "version": 8,
              "sources": {
                "world": {
                  "type": "geojson",
                  "data": `${this.$config.staticUri}/geo/world.json`,
                  "maxzoom": 6
                }
              },
              "glyphs": `${this.$config.staticUri}/font/{fontstack}/{range}.pbf`,
              "layers": [
                {
                  "id": "background",
                  "type": "background",
                  "paint": {
                    "background-color": "#aadafe"
                  }
                },
                {
                  id: "land",
                  type: "fill",
                  source: "world",
                  // "source-layer": "land",
                  paint: {
                    "fill-color": "#cbe5ca",
                  },
                },
                {
                  "id": "country-abbrev",
                  "type": "symbol",
                  "source": "world",
                  "maxzoom": 3,
                  "layout": {
                    "text-field": "{abbrev}",
                    "text-font": ["Open Sans Semibold"],
                    "text-transform": "uppercase",
                    "text-max-width": 20,
                    "text-size": {
                      "stops": [[3, 10], [4, 11], [5, 12], [6, 16]]
                    },
                    "text-letter-spacing": {
                      "stops": [[4, 0], [5, 1], [6, 2]]
                    },
                    "text-line-height": {
                      "stops": [[5, 1.2], [6, 2]]
                    }
                  },
                  "paint": {
                    "text-halo-color": "#fff",
                    "text-halo-width": 1
                  },
                },
                {
                  "id": "country-border",
                  "type": "line",
                  "source": "world",
                  "paint": {
                    "line-color": "#226688",
                    "line-opacity": 0.25,
                    "line-dasharray": [6, 2, 2, 2],
                    "line-width": 1.2
                  }
                },
                {
                  "id": "country-name",
                  "type": "symbol",
                  "minzoom": 3,
                  "source": "world",
                  "layout": {
                    "text-field": "{name}",
                    "text-font": ["Open Sans Semibold"],
                    "text-max-width": 20,
                    "text-size": {
                      "stops": [[3, 10], [4, 11], [5, 12], [6, 16]]
                    }
                  },
                  "paint": {
                    "text-halo-color": "#fff",
                    "text-halo-width": 1
                  },
                },
              ],
            },
            attributionControl: true,
            customAttribution: this.attribution,
            zoom: 0,
          };
        }

        this.url = 'https://api.maptiler.com/maps/' + s.style + '/{z}/{x}/{y}.png?key=' + mapKey;
        this.options = mapOptions;
        this.mapFont = mapFont;
      });
    },
    search() {
      const geojson = Photo.toGeoJson(this.photos);

      this.map.getSource("photos").setData(geojson);

      this.map.fitBounds(geojson.bbox, {
        maxZoom: 17,
        padding: 100,
        duration: this.settings.animate,
        essential: false,
        animate: this.settings.animate > 0
      });

      this.updateMarkers();
    },
    renderMap() {
      this.map = new maplibregl.Map(this.options);
      this.map.setLanguage(this.$config.values.settings.ui.language.split("-")[0]);

      this.map.addControl(new maplibregl.NavigationControl({showCompass: true}, 'top-right'));
      this.map.addControl(new maplibregl.FullscreenControl({container: document.querySelector('body')}));
      this.map.addControl(new maplibregl.GeolocateControl({
        positionOptions: {
          enableHighAccuracy: true
        },
        trackUserLocation: true
      }));

      this.map.on("load", () => this.onMapLoad());
    },
    updateMarkers() {
      let newMarkers = {};
      let features = this.map.querySourceFeatures("photos");

      for (let i = 0; i < features.length; i++) {
        let coords = features[i].geometry.coordinates;
        let props = features[i].properties;
        if (props.cluster) continue;
        let id = features[i].id;

        let marker = this.markers[id];
        let token = this.$config.previewToken();
        if (!marker) {
          let el = document.createElement('div');
          el.className = 'marker';
          el.title = props.Title;
          el.style.backgroundImage = `url(${this.$config.contentUri}/t/${props.Hash}/${token}/tile_50)`;
          el.style.width = '50px';
          el.style.height = '50px';

          el.addEventListener('click', () => this.openPhoto(props.UID));
          marker = this.markers[id] = new maplibregl.Marker({
            element: el
          }).setLngLat(coords);
        } else {
          marker.setLngLat(coords);
        }

        newMarkers[id] = marker;

        if (!this.markersOnScreen[id]) {
          marker.addTo(this.map);
        }
      }
      for (let id in this.markersOnScreen) {
        if (!newMarkers[id]) {
          this.markersOnScreen[id].remove();
        }
      }
      this.markersOnScreen = newMarkers;
    },
    onMapLoad() {
      this.map.addSource('photos', {
        type: 'geojson',
        data: null,
        cluster: true,
        clusterMaxZoom: 14, // Max zoom to cluster points on
        clusterRadius: 50 // Radius of each cluster when clustering points (defaults to 50)
      });

      this.map.addLayer({
        id: 'clusters',
        type: 'circle',
        source: 'photos',
        filter: ['has', 'point_count'],
        paint: {
          'circle-color': [
            'step',
            ['get', 'point_count'],
            '#2DC4B2',
            100,
            '#3BB3C3',
            750,
            '#669EC4'
          ],
          'circle-radius': [
            'step',
            ['get', 'point_count'],
            20,
            100,
            30,
            750,
            40
          ]
        }
      });

      this.map.addLayer({
        id: 'cluster-count',
        type: 'symbol',
        source: 'photos',
        filter: ['has', 'point_count'],
        layout: {
          'text-field': '{point_count_abbreviated}',
          'text-font': this.mapFont,
          'text-size': 13
        }
      });

      this.map.on('render', this.updateMarkers);

      this.map.on('click', 'clusters', (e) => {
        const features = this.map.queryRenderedFeatures(e.point, {
          layers: ['clusters']
        });
        const clusterId = features[0].properties.cluster_id;
        this.map.getSource('photos').getClusterExpansionZoom(
          clusterId,
          (err, zoom) => {
            if (err) return;

            this.map.easeTo({
              center: features[0].geometry.coordinates,
              zoom: zoom
            });
          }
        );
      });

      this.map.on('mouseenter', 'clusters', () => {
        this.map.getCanvas().style.cursor = 'pointer';
      });
      this.map.on('mouseleave', 'clusters', () => {
        this.map.getCanvas().style.cursor = '';
      });

      this.search();
    },
  },
};
</script>

