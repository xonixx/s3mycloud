<script setup lang="ts">
defineProps<{
  page: number
  totalPages: number
  changePage: (p: number) => {}
}>();
</script>

<template>
  <!--  <nav class="pagination" role="navigation" aria-label="pagination">-->
  <!--    <a class="pagination-previous is-disabled">Previous</a>
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
      </ul>-->
  <!--  <div class="columns is-centered">-->
  <!--    <div class="column is-2 field is-grouped">-->
  <div class="field is-grouped">
    <div class="control">
      <a class="button is-small" :disabled="true">
        <i class="fas fa-solid fa-angles-left"></i>
      </a>
    </div>
    <div class="control">
      <a class="button is-small">
        <i class="fas fa-solid fa-angle-left"></i>
      </a>
    </div>
    <div class="control" style="margin-right: 0">
      <input class="input is-small" type="text" placeholder="" style="width: 30px">
    </div>
    <div class="control">
      <input class="input is-small" type="text" disabled :value="'/ ' + totalPages" style="width: 40px; background-color: #fff; border-color: #fff">
<!--      <span class="input is-small" style="width: 40px; background-color: #fff; border:0">/ {{totalPages}}</span>-->
      <!--      / {{ totalPages }}-->
    </div>
    <div class="control">
      <a class="button is-small">
        <i class="fas fa-solid fa-angle-right"></i>
      </a>
    </div>
    <div class="control">
      <a class="button is-small">
        <i class="fas fa-solid fa-angles-right"></i>
      </a>
    </div>
  </div>
  <!--  </div>-->
  <!--  </nav>-->
</template>