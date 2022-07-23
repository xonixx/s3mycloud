import { ref } from "vue";
import { UploadableFile } from "./UploadableFile";

export default function () {
  const files = ref<UploadableFile[]>([]);

  function addFiles(newFiles: FileList) {
    let newUploadableFiles = [...newFiles]
      .map((file) => new UploadableFile(file))
      .filter((file) => !fileExists(file.id));
    files.value = files.value.concat(newUploadableFiles);
  }

  function fileExists(otherId:string):boolean {
    return files.value.some(({ id }) => id === otherId);
  }

  function removeFile(file) {
    const index = files.value.indexOf(file);

    if (index > -1) files.value.splice(index, 1);
  }

  return { files, addFiles, removeFile };
}
