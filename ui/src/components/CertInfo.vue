<template>
  <div v-if="certInfo">
    <button class="btn btn-primary mb-1" v-on:click="showCertInfo">
      See Certificate info
    </button>
    <table class="small table d-none text-light" ref="certInfoTable">
      <tr v-for="(value, key) in certInfo" :key="key">
        <td>{{ key }}</td>
        <td>{{ value }}</td>
      </tr>
    </table>
  </div>
  <div v-else class="small">
    <p>No certificate information</p>
    <ul>
      <li>You didn't provide a certificate for this service</li>
      <li>
        If you use a wildcard certificate, make a call to the server to fill the
        cache
      </li>
      <li>A self-signed certificate will be generated in the others cases</li>
    </ul>
  </div>
</template>
<script lang="ts">
import Certificate from "@/services/certificate";
import { Options, Vue } from "vue-class-component";

@Options({
  name: "CertInfo",
  props: {
    certInfo: {
      type: Object,
      required: false,
    },
  },
})
export default class CertInfoComponent extends Vue {
  certInfo!: Certificate | null;
  showCertInfo() {
    (this.$refs.certInfoTable as HTMLElement).classList.toggle("d-none");
  }
}
</script>
