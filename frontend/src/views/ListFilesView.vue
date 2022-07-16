<script setup lang="ts">

import {ref} from "vue";

const searchResults = ref([])

for (let i = 0; i < 20; i++) {
  searchResults.value.push(
      {name: "Report for boss.xlsx", size: 50000, date: "15 Mar 2016", tags: ["document", "work"]},
      {name: "Sing Now.mp3", size: 2_500_000, date: "17 Apr 2019", tags: ["music", "pop"]},
      {name: "Test.txt", size: 100, date: "1 Jan 2008", tags: []},
      {name: "CV (John Doe).pdf", size: 123_456, date: "9 Mar 2020", tags: ["work"]},
      {name: "Some veeeeeeery loooooooooooong naaaaaaaaame.ext", size: 0, date: "31 Jan 2010", tags: ["test"]},
  )
}


</script>

<template>
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
  <div>
    <div v-for="r in searchResults" class="res-item">
      <span class="name">{{ r.name }}</span>
      <span v-for="t in r.tags" class="tag is-primary is-light">{{ t }}</span>
      <div style="margin-top: -5px">
        <span class="date">{{ r.date }}</span>
        <span class="size">{{ r.size }}</span>
      </div>
    </div>
  </div>
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
