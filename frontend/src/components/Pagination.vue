<script setup lang="ts">
const props = defineProps<{
  page: number;
  pageSize: number;
  totalRecords: number;
  changePage: (p: number) => void;
  changePageSize: (ps: number) => void;
}>();

const pageSizes = [10, 50, 100];

function totalPages(): number {
  console.info("tr", props.totalRecords, "ps", props.pageSize);
  return Math.ceil(props.totalRecords / props.pageSize);
}

function myChangePage(p: number): void {
  let tot: number;
  if (p < 0) p = 0;
  else if (p >= (tot = totalPages())) p = tot - 1;
  props.changePage(p);
}

function isFirstPage(): boolean {
  return props.page == 0;
}

function isLastPage(): boolean {
  return props.page == totalPages() - 1;
}

function isNavigationVisible(): boolean {
  return props.totalRecords > props.pageSize;
}

function isPageSizeVisible(): boolean {
  return props.totalRecords > pageSizes[0];
}
</script>

<template>
  <!--  <div class="columns is-centered">-->
  <!--    <div class="column is-2 field is-grouped">-->
  <div class="field is-grouped">
    <template v-if="isNavigationVisible()">
      <div class="control">
        <button
          class="button is-small"
          :disabled="isFirstPage()"
          @click="myChangePage(0)"
        >
          <i class="fas fa-solid fa-angles-left"></i>
        </button>
      </div>
      <div class="control">
        <button
          class="button is-small"
          :disabled="isFirstPage()"
          @click="myChangePage(page - 1)"
        >
          <i class="fas fa-solid fa-angle-left"></i>
        </button>
      </div>
      <div class="control" style="margin-right: 0">
        <input
          class="input is-small"
          type="text"
          placeholder=""
          :value="page + 1"
          @keydown.enter="
            myChangePage((parseInt($event.target.value) || 1) - 1)
          "
          style="width: 50px"
        />
      </div>
      <div class="control">
        <!--      <input class="input is-small" type="text" disabled :value="'/ ' + totalPages()"-->
        <!--             style="width: 40px; background-color: #fff; border-color: #fff">-->
        <span
          class="input is-small"
          style="width: 50px; background-color: #fff; border: 0"
          >/ {{ totalPages() }}</span
        >
        <!--      / {{ totalPages }}-->
      </div>
      <div class="control">
        <button
          class="button is-small"
          :disabled="isLastPage()"
          @click="myChangePage(page + 1)"
        >
          <i class="fas fa-solid fa-angle-right"></i>
        </button>
      </div>
      <div class="control" style="margin-right: 50px">
        <button
          class="button is-small"
          :disabled="isLastPage()"
          @click="myChangePage(totalPages() - 1)"
        >
          <i class="fas fa-solid fa-angles-right"></i>
        </button>
      </div>
    </template>
    <template v-if="isPageSizeVisible()">
      <div class="control">
        <div class="select is-small">
          <select
            :value="pageSize"
            @change="changePageSize($event.target.value)"
          >
            <option v-for="ps in pageSizes" :value="ps">{{ ps }}</option>
          </select>
        </div>
      </div>
      <div class="control">per page</div>
    </template>
  </div>
  <!--  </div>-->
</template>
