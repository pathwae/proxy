<template>
  <nav>
    <img alt="Pathwae logo" src="../assets/logo.svg" class="logo" />
    <ul class="nav">
      <li class="nav-item">
        <router-link to="/"><i class="bi bi-activity"></i> Stats</router-link>
      </li>
      <li class="nav-item">
        <router-link to="/about"> <i class="bi bi-info"></i> About</router-link>
      </li>
    </ul>
  </nav>
</template>
<style lang="scss">
@import "@/assets/global.scss";
@import "bootstrap-icons/font/bootstrap-icons.css";
nav {
  background-color: $gray-700;
  @extend .shadow;
  li {
    a {
      display: inline-block;
      @extend .px-2;
      @extend .py-1;
      color: $white;
      &:hover {
        color: $gray-300;
      }
      &.router-link-exact-active {
        background-color: $gray-600;
        color: $white;
      }
    }
  }
}
</style>
