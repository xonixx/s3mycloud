<script setup lang="ts">

import {ref} from "vue";

const searchResults = ref([])

for (let i = 0; i < 2; i++) {
  searchResults.value.push(
      {name: "Report for boss.xlsx", size: 50000, date: "15 Mar 2016", tags: ["document", "work"]},
      {name: "Sing Now.mp3", size: 2_500_000, date: "17 Apr 2019", tags: ["music", "pop"]},
      {name: "Test.txt", size: 100, date: "1 Jan 2008", tags: []},
      {name: "CV (John Doe).pdf", size: 123_456, date: "9 Mar 2020", tags: ["work"]},
      {name: "Some veeeeeeery loooooooooooong naaaaaaaaame.ext", size: 0, date: "31 Jan 2010", tags: ["test"]},
  )
}

/**
 * Format bytes as human-readable text.
 *
 * @param bytes Number of bytes.
 * @param si True to use metric (SI) units, aka powers of 1000. False to use
 *           binary (IEC), aka powers of 1024.
 * @param dp Number of decimal places to display.
 *
 * @return Formatted string.
 */
function humanFileSize(bytes: number, si=false, dp=1) {
  const thresh = si ? 1000 : 1024;

  if (Math.abs(bytes) < thresh) {
    return bytes + ' B';
  }

  const units = si
      ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
      : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let u = -1;
  const r = 10**dp;

  do {
    bytes /= thresh;
    ++u;
  } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


  return bytes.toFixed(dp) + ' ' + units[u];
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
        <span class="size">{{ humanFileSize(r.size,true) }}</span>
      </div>
    </div>
  </div>

  <nav class="pagination" role="navigation" aria-label="pagination">
    <a class="pagination-previous">Previous</a>
    <a class="pagination-next">Next page</a>
    <ul class="pagination-list">
      <li>
        <a class="pagination-link" aria-label="Goto page 1">1</a>
      </li>
      <li>
        <span class="pagination-ellipsis">&hellip;</span>
      </li>
      <li>
        <a class="pagination-link" aria-label="Goto page 45">45</a>
      </li>
      <li>
        <a class="pagination-link is-current" aria-label="Page 46" aria-current="page">46</a>
      </li>
      <li>
        <a class="pagination-link" aria-label="Goto page 47">47</a>
      </li>
      <li>
        <span class="pagination-ellipsis">&hellip;</span>
      </li>
      <li>
        <a class="pagination-link" aria-label="Goto page 86">86</a>
      </li>
    </ul>
  </nav>
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
