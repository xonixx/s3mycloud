<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
import HelloWorld from "./components/HelloWorld.vue";
</script>

<template>
  <section class="section">
    <div class="container">
      <h1 class="title">
        S3MyCloud
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

<style scoped>

</style>
