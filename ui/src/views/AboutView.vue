<template>
  <div class="about container mt-4">
    <h1>What is Pathwae?</h1>
    <p>
      Pathwae is an opensource and free web http/https proxy build with
      <a href="https://go.dev" target="_blank">Go</a>.
    </p>
    <p>
      The web UI is built with
      <a href="https://vuejs.org" target="_blank">Vue</a> using
      <a href="https://www.typescriptlang.org/" target="_blank">Typescript</a>.
    </p>
    <blockquote>It's a Reverse Proxy, nothing more, nothing less.</blockquote>
    <p>
      It's goal is to make it possible to use it with Podman and Docker of
      development environment. You only need to set a
      <code>CONFIG</code> environment variable or <code>CONFIG_FILE</code> path
      to a YAML or JSON content where you define routes and enpoints.
    </p>
    <p>
      Pathwae automatically finds SSL/TLS certificates that are written in
      <code>/certs</code> directory, you don't need to specify wich one to use
      for specific route. If you set a wildcard certificate, so Pathwae will
      manage how to use it. If no certificate is provided for a route, Pathwae
      generates a temporary certificate.
    </p>
    <p>
      Exmpple using docker-compose file to server a
      <a href="https://ghost.org/" target="_blank">Ghost blog</a>:
    </p>
    <pre>
    <code>
version: '3'
services:
    blog:
      image: ghost
      environment:
        url: "https://blog.proj.localhost"

    proxy:
      image: pathwae/pathwae
      environment:
        CONFIG: |
          blog.proj.localhost:
            to: http://blog:2368
            force_ssl: true
      ports:
        # Pathwae is a HTTP(s) server
        - "80:80"
        - "443:443"
        # api / webui
        - "8080:8080"
    </code>
    </pre>
    <p>
      You can also use <code>CONFIG_FILE</code> environment variable to specify
      a path to a YAML or JSON file.
    </p>
    <p>
      The configuration is a list of routes maps where the key is the
      "entrypoint". Each route is a map with the following keys:
    </p>
    <ul>
      <li>
        <code>to</code>: the endpoint to proxy to. The scheme is important
        (default is http)
      </li>
      <li>
        <code>force_ssl</code>: if true, the route will be redirected to HTTPS
      </li>
      <li>
        <code>enabled</code>: if false, the route will be ignored (default is
        true)
      </li>
    </ul>
  </div>
</template>
<script lang="ts">
import hljs from "highlight.js";
import "highlight.js/scss/atom-one-dark.scss";
import { Vue } from "vue-class-component";
export default class AboutView extends Vue {
  mounted() {
    document.querySelectorAll("pre code").forEach((block) => {
      hljs.highlightBlock(block as HTMLElement);
    });
  }
}
</script>
