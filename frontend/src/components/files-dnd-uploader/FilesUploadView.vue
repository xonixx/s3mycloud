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
      <ul class="image-list" v-show="files.length">
        <FilePreview
          v-for="file of files"
          :key="file.id"
          :file="file"
          tag="li"
          @remove="removeFile"
        />
      </ul>
    </DropZone>
    <button @click.prevent="uploadFiles(files)" class="upload-button">
      Upload
    </button>
  </div>
</template>

<script setup lang="ts">
// Components
import DropZone from "./DropZone.vue";
import FilePreview from "./FilePreview.vue";

// File Management
import useFileList from "./file-list";

const { files, addFiles, removeFile } = useFileList();

function onInputChange(e:Event & { target: HTMLInputElement }) {
  addFiles(e.target.files);
  e.target.value = null; // reset so that selecting the same file again will still cause it to fire this change
}

// Uploader
import createUploader from "./file-uploader";

const { uploadFiles } = createUploader("YOUR URL HERE");
</script>

<!--<style lang="css">
html {
  height: 100%;
  width: 100%;
  background-color: #b6d6f5;

  /* Overlapping Stripes Background, based off https://codepen.io/MinzCode/pen/Exgpqpe */
  &#45;&#45;color1: rgba(55, 165, 255, 0.35);
  &#45;&#45;color2: rgba(96, 181, 250, 0.35);
  &#45;&#45;rotation: 135deg;
  &#45;&#45;size: 10px;
  background-blend-mode: multiply;
  background-image: repeating-linear-gradient(
      var(&#45;&#45;rotation),
      var(&#45;&#45;color1) calc(var(&#45;&#45;size) * 0),
      var(&#45;&#45;color1) calc(var(&#45;&#45;size) * 9),
      var(&#45;&#45;color2) calc(var(&#45;&#45;size) * 9),
      var(&#45;&#45;color2) calc(var(&#45;&#45;size) * 12),
      var(&#45;&#45;color1) calc(var(&#45;&#45;size) * 12)
    ),
    repeating-linear-gradient(
      calc(var(&#45;&#45;rotation) + 90deg),
      var(&#45;&#45;color1) calc(var(&#45;&#45;size) * 0),
      var(&#45;&#45;color1) calc(var(&#45;&#45;size) * 9),
      var(&#45;&#45;color2) calc(var(&#45;&#45;size) * 9),
      var(&#45;&#45;color2) calc(var(&#45;&#45;size) * 12),
      var(&#45;&#45;color1) calc(var(&#45;&#45;size) * 12)
    );
}

body {
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>-->

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
}

.drop-area {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding: 50px;
  background: #ffffff55;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
  transition: .2s ease;

  &[data-active=true] {
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

  input[type=file]:not(:focus-visible) {
    /*Visually Hidden Styles taken from Bootstrap 5*/
    position:absolute !important;
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
