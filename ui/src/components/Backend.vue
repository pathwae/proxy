<template>
  <div>
    <div :class="status" class="shadow-lg p-3 m-1 rounded">
      <h2>
        <!-- from address -->
        <div class="d-flex justify-content-space-between">
          <div class="input-group input-group-sm">
            <i class="mr-2 bi" :class="state"></i>
            <a
              target="_blank"
              :href="scheme + '://' + backend.name"
              ref="backend-in"
              class="small"
              >{{ backend.name }}</a
            >
          </div>
          <div class="d-flex">
            <button
              class="btn"
              :class="{
                'bi-pause-fill': backend.enabled,
                'bi-play-fill': !backend.enabled,
              }"
              title="pause/unpause"
              v-on:click="toggleState"
              ref="actionButton"
            ></button>
          </div>
        </div>
      </h2>
      <!-- To backend -->
      <div class="d-flex justify-content-space-between">
        <div class="input-group input-group-sm">
          <span class="mr-1">To:</span>
          <span ref="backend-to" :class="{ 'd-none': toEditionInProgress }">{{
            backend.to
          }}</span>
          <input
            type="text"
            class="form-control"
            :class="{ 'd-none': !toEditionInProgress }"
            ref="edit-backend-to"
            v-model="backend.to"
            v-on:keyup.enter="updateCurrentBackend"
          />
        </div>
        <div>
          <button
            class="btn bi bi-pencil"
            v-on:click="editBackendTo"
            title="Edit the backend address"
          ></button>
        </div>
      </div>
      <p>
        <strong class="mr-2">Total Hits {{ hits }}</strong>
        <span
          v-for="(count, method) in perMethods"
          :key="method"
          class="small mr-2"
        >
          {{ method }}: {{ count }}
        </span>
      </p>
      <CertInfo :certInfo="certInfo" />
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Options } from "vue-class-component";
import Backend from "@/services/backend";
import Api from "@/services/api";
import Hit from "@/services/hit";
import Certificate from "@/services/certificate";
import Store from "@/services/store";
import CertInfo from "@/components/CertInfo.vue";

const OK_CHAR = "success bi-check-circle-fill";
const KO_CHAR = "error bi-exclamation-circle";
const UNKNOWN_CHAR = "bi-question-circle";

@Options({
  name: "BackendComponent",
  components: {
    CertInfo,
  },
  props: {
    backend: {
      type: Object,
      required: true,
    },
  },
})
export default class BackendComponent extends Vue {
  backend!: Backend;

  certInfo: Certificate | null = null;
  scheme = "http";

  // view elements
  state = UNKNOWN_CHAR;
  stats = new Array<Hit>();
  hits = 0;
  status = "";
  perMethods = {} as { [key: string]: number };

  // utlils
  api = new Api();
  sse: EventSource | null = null;
  store = Store;

  toEditionInProgress = false;

  /**
   * When the component is mounted...
   */
  mounted() {
    this.initialize();
    this.startListening();
  }

  /**
   * Stop listening to the backend Server Event.
   */
  unmounted() {
    if (this.sse !== null) {
      this.sse.close();
    }
  }

  editBackendTo() {
    this.toEditionInProgress = !this.toEditionInProgress;
  }

  /**
   * Toggle the state of the backend.
   */
  toggleState() {
    this.backend.enabled = !this.backend.enabled;
    this.updateCurrentBackend();
  }

  /**
   * Update the backend input address.
   */
  async updateCurrentBackend() {
    this.toEditionInProgress = false;
    this.api.setBackend(this.backend.name, this.backend).then(
      () => {
        console.log("Done");
      },
      (err) => {
        console.error(err);
      }
    );
  }

  /**
   * Initialize the component.
   */
  initialize() {
    this.scheme = this.backend.force_ssl ? "https" : "http";
    this.api.getCert(this.backend.name).then(
      (cert) => {
        this.certInfo = cert;
      },
      () => {
        this.certInfo = null;
      }
    );

    // Initialize the stats
    this.api.getBackendStats(this.backend.name).then((stats: Array<Hit>) => {
      if (stats === null) {
        return;
      }
      // count the number of hits
      this.stats = stats;
      this.hits = stats.length;

      for (const stat of stats) {
        if (this.perMethods[stat.method] === undefined) {
          this.perMethods[stat.method] = 0;
        }
        this.perMethods[stat.method]++;
        this.store.incrementHits();
      }
    });
  }

  /**
   * Start listening to the backend Server Event.
   */
  startListening() {
    if (this.sse !== null) {
      this.sse.close();
    }

    // Update the status and the state
    const sse = new EventSource(
      `${this.api.BASE_URL}/sse/status/${this.backend.name}`
    );
    this.sse = sse;
    // when status of the backend changes
    sse.addEventListener("status", (e: Event) => {
      // if e has no property named data, return
      const data = JSON.parse((e as MessageEvent).data);
      this.state = data.Stat ? OK_CHAR : KO_CHAR;
      if (this.backend.enabled) {
        this.status = data.Stat ? "no-problem" : "bg-danger";
      } else {
        this.status = "bg-warning text-dark opacity-50";
      }
    });

    // got stat (hits)
    sse.addEventListener("stat", (e: Event) => {
      // if e has no property named data, return
      const stat = JSON.parse((e as MessageEvent).data);
      if (this.perMethods[stat.method] === undefined) {
        this.perMethods[stat.method] = 0;
      }
      this.perMethods[stat.method]++;
      this.stats.push(stat);
      this.hits++;
      this.store.incrementHits();
    });

    // got cert
    sse.addEventListener("cert", (e: Event) => {
      // if e has no property named data, return
      const cert = JSON.parse((e as MessageEvent).data);
      this.certInfo = cert;
    });
  }
}
</script>
<style lang="scss">
@import "@/assets/global.scss";

.bi {
  &.success {
    color: $green;
  }
  &.error {
    color: $white;
  }
}

.no-problem {
  background-color: var($secondary);
}
.opacity-50 {
  opacity: 0.5;
}
</style>
