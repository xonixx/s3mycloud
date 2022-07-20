<script setup lang="ts">
import {ref} from "vue";

const props = defineProps<{
  page: number
  pageSize: number
  totalRecords: number
  changePage: (p: number) => void
}>();

const pageRef = ref()

function totalPages(): number {
  console.info("tr", props.totalRecords,"ps",props.pageSize)
  return props.totalRecords / props.pageSize;
}

function myChangePage(p: number): void {
  let tot:number;
  if (p < 0) p = 0;
  else if (p >= (tot = totalPages()))
    p = tot - 1
  props.changePage(p)
}

</script>

<template>
  PR {{ pageRef }} tr {{ totalRecords }}
  <!--  <div class="columns is-centered">-->
  <!--    <div class="column is-2 field is-grouped">-->
  <div class="field is-grouped">
    <div class="control">
      <button class="button is-small" :disabled="true">
        <i class="fas fa-solid fa-angles-left"></i>
      </button>
    </div>
    <div class="control">
      <button class="button is-small" @click="myChangePage(page-1)">
        <i class="fas fa-solid fa-angle-left"></i>
      </button>
    </div>
    <div class="control" style="margin-right: 0">
      <input class="input is-small" type="text" placeholder=""
             :value="page+1"
             @keyup.enter="myChangePage((parseInt(pageRef.value) || 1)-1)"
             ref="pageRef"
             style="width: 50px">
    </div>
    <div class="control">
<!--      <input class="input is-small" type="text" disabled :value="'/ ' + totalPages()"-->
<!--             style="width: 40px; background-color: #fff; border-color: #fff">-->
            <span class="input is-small" style="width: 50px; background-color: #fff; border:0">/ {{totalPages()}}</span>
      <!--      / {{ totalPages }}-->
    </div>
    <div class="control">
      <button class="button is-small" @click="myChangePage(page+1)">
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