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

/*for (let i = 0; i < 2; i++) {
  searchResults.value.push(
      {name: "Report for boss.xlsx", size: 50000, date: "15 Mar 2016", tags: ["document", "work"]},
      {name: "Sing Now.mp3", size: 2_500_000, date: "17 Apr 2019", tags: ["music", "pop"]},
      {name: "Test.txt", size: 100, date: "1 Jan 2008", tags: []},
      {name: "CV (John Doe).pdf", size: 123_456, date: "9 Mar 2020", tags: ["work"]},
      {name: "Some veeeeeeery loooooooooooong naaaaaaaaame.ext", size: 0, date: "31 Jan 2010", tags: ["test"]},
  )
}*/

async function load() {
  const response = await axios.get(`http://127.0.0.1:8080/api/file?page=${params.page}&pageSize=${params.pageSize}`);
  searchResults.value = response.data.page.map((e) => ({...e, date: moment(e.uploaded).format("D MMM YYYY")}));
  totalRecords.value = response.data.total;
}

onMounted(load)

function changePage(p: number) {
  console.info("Changing page to ", p)
  router.push({ path: '/', query: { ...router.currentRoute.value.query, page: p } })
  setTimeout(load,10)
}

</script>

<template>
  PS:{{ pageSize }} P:{{ page }}
  <!--  <main>-->
  <div class="field is-grouped">
    <div class="control is-expanded">
      <input class="input" type="text" placeholder="Enter search query...">
    </div>
    <div class="control">
      <a class="button is-info">
        Search
      </a>
    </div>
  </div>

  <Pagination :page="page" :page-size="pageSize" :total-records="totalRecords" :change-page="changePage"/>

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

<!--  <Pagination/>-->
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
