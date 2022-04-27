<template>
  <div class="home">
    <div class="p-4">
      <div class="container">
        <h2>Quick view</h2>
        <div class="row">
          <InfoBlock
            title="Backends"
            :message="store.backends.length.toString()"
            class="info-block bg-primary"
            icon="boxes"
          />
          <InfoBlock
            title="Memory"
            :message="serverMemory"
            class="info-block bg-success"
            icon="memory"
          />
          <InfoBlock
            title="Requests"
            :message="store.hits.toString()"
            class="info-block bg-warning"
            icon="server"
          />
          <InfoBlock
            title="Version"
            :message="version"
            class="info-block bg-danger"
            icon="info"
          />
        </div>
      </div>
    </div>
    <div class="bg-dark p-4">
      <div class="container">
        <h2>Routes</h2>
        <div class="row justify-content-center">
          <BackendComponent
            v-for="(b, idx) in store.backends"
            :key="idx"
            :backend="b"
            class="col-md-6 p-0"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Api from "@/services/api";
import BackendComponent from "@/components/Backend.vue";
import Store from "@/services/store";
import InfoBlock from "@/components/InfoBlock.vue";
import Backend from "@/services/backend";

@Options({
  components: { BackendComponent, InfoBlock },
})
export default class HomeView extends Vue {
  serverMemory = "0";
  version = "unknown";
  store = Store;
  interval!: number;
  sse: EventSource | null = null;

  api = new Api();

  mounted() {
    this.store.hits = 0; // the counter will be reset when all servers are contacted

    this.api.getBackends().then((backends) => {
      const bs = new Array<Backend>();
      for (const name in backends) {
        const backend = backends[name];
        backend.name = name;
        bs.push(backend);
      }
      this.store.backends = bs;
    });

    this.api.getServerVersion().then((version) => {
      this.version = version;
    });

    const sse = new EventSource(`${this.api.BASE_URL}/sse/global`);

    // got memory usage from proxy
    sse.addEventListener("memory", (mem) => {
      const memory = mem as MessageEvent;
      this.serverMemory = (memory.data / 1024 / 1024).toFixed(2) + " Mb";
    });

    // backend has changed, update the list
    sse.addEventListener("changes", (edited) => {
      const editedEvent = edited as MessageEvent;
      let backend = this.store.backends.find(
        (b) => b.name === editedEvent.data.name
      );
      if (backend) {
        backend = editedEvent.data.backend;
      }
    });

    this.sse = sse;
  }

  unmounted() {
    if (this.sse !== null) {
      this.sse.close();
    }
  }
}
</script>
<style lang="scss">
@import "@/assets/global.scss";

.info-block {
  @extend .col-md;
  @extend .shadow;
  @extend .rounded;
  @extend .m-md-2;
  @extend .p-md-2;
  @extend .mb-1;
  @extend .text-center;
}
</style>
