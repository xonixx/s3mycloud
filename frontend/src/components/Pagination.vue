<script setup lang="ts">
import {ref} from "vue";

defineProps<{
  page: number
  totalPages: number
  changePage: (p: number) => void
}>();

const pageRef = ref()

</script>

<template>
  PR {{ pageRef }}
  <!--  <div class="columns is-centered">-->
  <!--    <div class="column is-2 field is-grouped">-->
  <div class="field is-grouped">
    <div class="control">
      <button class="button is-small" :disabled="true">
        <i class="fas fa-solid fa-angles-left"></i>
      </button>
    </div>
    <div class="control">
      <button class="button is-small" @click="changePage(page-1)">
        <i class="fas fa-solid fa-angle-left"></i>
      </button>
    </div>
    <div class="control" style="margin-right: 0">
      <input class="input is-small" type="text" placeholder=""
             :value="page+1"
             @keyup.enter="changePage((parseInt(pageRef.value) || 1)-1)"
             ref="pageRef"
             style="width: 30px">
    </div>
    <div class="control">
      <input class="input is-small" type="text" disabled :value="'/ ' + totalPages"
             style="width: 40px; background-color: #fff; border-color: #fff">
<!--      <span class="input is-small" style="width: 40px; background-color: #fff; border:0">/ {{totalPages}}</span>-->
      <!--      / {{ totalPages }}-->
    </div>
    <div class="control">
      <button class="button is-small" @click="changePage(page+1)">
        <i class="fas fa-solid fa-angle-right"></i>
      </button>
    </div>
    <div class="control">
      <button class="button is-small">
        <i class="fas fa-solid fa-angles-right"></i>
      </button>
    </div>
  </div>
  <!--  </div>-->
</template>