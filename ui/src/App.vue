<template>
  <div>
    <NavBar class="p-3 navbar sticky-top" />
    <div>
      <router-view />
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Options } from "vue-class-component";
import NavBar from "@/components/NavBar.vue";

@Options({
  components: {
    NavBar,
  },
})
export default class App extends Vue {}
</script>
