<script setup lang="ts">
import {humanFileSize} from "@/assets/lib";
import Pagination from "@/components/Pagination.vue";

import moment from 'moment';
import {onMounted, ref} from "vue";
import axios from "axios";
import router from "@/router";

const params = withDefaults(defineProps<{
  q: string;
  page: number;
  pageSize: number;
}>(), {
  q: '',
  page: 0,
  pageSize: 10
});

const totalRecords = ref<number>();
const searchResults = ref([])

async function load() {
  const response = await axios.get(`http://127.0.0.1:8080/api/file?page=${params.page}&pageSize=${params.pageSize}&name=${params.q}`);
  searchResults.value = response.data.page.map((e) => ({...e, date: moment(e.uploaded).format("D MMM YYYY")}));
  totalRecords.value = response.data.total;
}

onMounted(load)

function changePage(p: number) {
  console.info("Changing page to ", p)
  router.push({path: '/', query: {...router.currentRoute.value.query, page: p}})
  setTimeout(load, 10)
}

function changePageSize(ps: number) {
  console.info("Changing page size to ", ps)
  router.push({path: '/', query: {...router.currentRoute.value.query, pageSize: ps}})
  setTimeout(load, 10)
}

const qRef = ref()

function changeQuery(q: string) {
  console.info("Changing query to ", q)
  router.push({path: '/', query: {...router.currentRoute.value.query, q, page: 0}})
  setTimeout(load, 10)
}

</script>

<template>
  <!--  <main>-->
  <div class="field is-grouped">
    <div class="control is-expanded">
      <input ref="qRef" :value="q" @keydown.enter="changeQuery(qRef.value)"
             class="input" type="text" placeholder="Enter search query...">
    </div>
    <div class="control">
      <a class="button is-info">
        Search
      </a>
    </div>
  </div>

  <Pagination
      :page="page"
      :page-size="pageSize"
      :total-records="totalRecords"
      :change-page="changePage"
      :change-page-size="changePageSize"
  />

  <div>
    <div v-for="r in searchResults" class="res-item">
      <a :href="`http://127.0.0.1:8080/api/file/${r.id}/dl`" class="name">{{ r.name }}</a>
      <span v-for="t in r.tags" class="tag is-primary is-light">{{ t }}</span>
      <div style="margin-top: -5px">
        <span class="date">{{ r.date }}</span>
        <span class="size">{{ humanFileSize(r.size, true) }}</span>
      </div>
    </div>
  </div>

  <Pagination
      :page="page"
      :page-size="pageSize"
      :total-records="totalRecords"
      :change-page="changePage"
      :change-page-size="changePageSize"
  />

  <!--  </main>-->
</template>

<style scoped>
.res-item {
  /*border: dotted 1px gray;*/
  margin-bottom: 10px;
}

.res-item span {
  margin: 5px;
}

.name {
  font-weight: bold;
}

.size {
  /*width: 90px;*/
  /*display: inline-block;*/
}

.res-item .date {
  margin-right: 30px;
  /*font-size: 0.9em;*/
  color: #888;
  /*font-style: italic;*/
}
</style>
