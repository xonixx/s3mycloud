<template>
  <div id="app">
    <DropZone
      class="drop-area"
      @files-dropped="addFiles"
      #default="{ dropZoneActive }"
    >
      <label for="file-input">
        <span v-if="dropZoneActive">
          <span>Drop Them Here</span>
          <span class="smaller">to add them</span>
        </span>
        <span v-else>
          <span>Drag Your Files Here</span>
          <span class="smaller">
            or <strong><em>click here</em></strong> to select files
          </span>
        </span>

        <input type="file" id="file-input" multiple @change="onInputChange" />
      </label>
      <!--      <ul class="image-list" v-show="files.length">
              <FilePreview
                v-for="file of files"
                :key="file.id"
                :file="file"
                tag="li"
                @remove="removeFile"
              />
            </ul>-->

      <template v-if="files.length">
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">Tags</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <input
                  v-model="tags"
                  class="input"
                  type="text"
                  placeholder="Tags to assign (comma-separated)"
                />
              </div>
            </div>
          </div>
        </div>

        <table class="table is-fullwidth files-table">
          <tr>
            <th>Name</th>
            <th>Size</th>
            <th>Status</th>
            <th></th>
          </tr>
          <tr v-for="file in files" :key="file.id">
            <td>{{ file.file.name }}</td>
            <td>{{ humanFileSize(file.file.size, true) }}</td>
            <td>
              <span v-if="file.status === 'loading'">In Progress</span>
              <span v-else-if="file.status === true">Uploaded</span>
              <span v-else-if="file.status === false">Error</span>
            </td>
            <td>
              <button @click="removeFile(file)">&times;</button>
            </td>
          </tr>
        </table>
      </template>
    </DropZone>
    <button
      @click.prevent="uploadFiles(files, splitTags(tags))"
      class="button is-primary upload-button"
      :disabled="!files.length || !tags"
    >
      Upload
    </button>
  </div>
</template>

<script setup lang="ts">
import { humanFileSize } from "@/assets/lib";
// Components
import DropZone from "./DropZone.vue";
import FilePreview from "./FilePreview.vue";

// File Management
import useFileList from "./file-list";

const { files, addFiles, removeFile } = useFileList();

const tags = ref<string>();

function splitTags(tags: string): string[] {
  tags = (tags || "").trim();
  if (!tags.length) {
    return [];
  }
  return tags.split(/\s*,\s*/);
}

function onInputChange(e: Event & { target: HTMLInputElement }) {
  addFiles(e.target.files);
  e.target.value = null; // reset so that selecting the same file again will still cause it to fire this change
}

// Uploader
import createUploader from "./file-uploader";
import { ref } from "vue";

const { uploadFiles } = createUploader();
</script>

<style scoped lang="scss">
#app {
  font-family: Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin: 0 auto;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;

  .files-table {
    text-align: left;

    td:last-child {
      text-align: right;
    }
  }
}

.drop-area {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding: 50px;
  background: #ffffff55;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
  transition: 0.2s ease;

  &[data-active="true"] {
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
    background: #ffffffcc;
  }
}

label {
  font-size: 36px;
  cursor: pointer;
  display: block;

  span {
    display: block;
  }

  input[type="file"]:not(:focus-visible) {
    /*Visually Hidden Styles taken from Bootstrap 5*/
    position: absolute !important;
    width: 1px !important;
    height: 1px !important;
    padding: 0 !important;
    margin: -1px !important;
    overflow: hidden !important;
    clip: rect(0, 0, 0, 0) !important;
    white-space: nowrap !important;
    border: 0 !important;
  }

  .smaller {
    font-size: 16px;
  }
}

.image-list {
  display: flex;
  list-style: none;
  flex-wrap: wrap;
  padding: 0;
}

.upload-button {
  display: block;
  appearance: none;
  border: 0;
  border-radius: 50px;
  padding: 0.75rem 3rem;
  margin: 1rem auto;
  font-size: 1.25rem;
  font-weight: bold;
  background: #369;
  color: #fff;
  text-transform: uppercase;
}

button {
  cursor: pointer;
}
</style>
