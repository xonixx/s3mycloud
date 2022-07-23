<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
import HelloWorld from "./components/HelloWorld.vue";
import router from "@/router";

function isActive(routeName: string): boolean {
  return routeName === router.currentRoute.value.name;
}
</script>

<template>
  <section class="section">
    <div class="container">
      <h1 class="title">
        S3MyCloud
        <span class="title is-4">
          <a :class="{ active: isActive('list') }" href="/">List</a>
          <a :class="{ active: isActive('upload') }" href="/upload">Upload</a>
        </span>
      </h1>
      <RouterView />
    </div>
  </section>
  <!--  <header>-->
  <!--    <div class="wrapper">-->
  <!--      <HelloWorld msg="You did it!" />-->

  <!--      <nav>-->
  <!--        <RouterLink to="/">Home</RouterLink>-->
  <!--        <RouterLink to="/about">About</RouterLink>-->
  <!--      </nav>-->
  <!--    </div>-->
  <!--  </header>-->
</template>

<style scoped lang="scss">
a {
  margin: 0 5px;

  &.active {
    background-color: #f5f5f5;
  }
}
</style>
