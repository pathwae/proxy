<template>
  <div>
    <i class="display-4 bi" :class="icon ? 'bi-' + icon : 'bi-info'"></i>
    <h3>{{ title }}</h3>
    <p class="lead">{{ message }}</p>
  </div>
</template>
<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  name: "InfoBlock",
  props: {
    title: {
      type: String,
      required: true,
    },
    message: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      required: false,
    },
  },
})
export default class InfoBlockComponent extends Vue {
  title!: string;
  message!: string;
  icon!: string;
}
</script>
