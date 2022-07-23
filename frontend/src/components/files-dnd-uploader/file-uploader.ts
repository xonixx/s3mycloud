import type { UploadableFile } from "./UploadableFile";
import { BACKEND_BASE_URL } from "@/urls";
import axios from "axios";

export async function uploadFile(file: UploadableFile) {
  // set up the request data
  let formData = new FormData();
  formData.append("file", file.file);

  // track status and upload file
  file.status = "loading";

  const response1: { data: { id: string; uploadUrl: string } } =
    await axios.post(BACKEND_BASE_URL + "/api/file/upload", {
      name: file.file.name,
      size: file.file.size,
      tags: ["TODO"],
    });

  // TODO handle error

  const response2 = await axios.put(response1.data.uploadUrl, formData);

  // TODO handle error

  // change status to indicate the success of the upload request
  // file.status = response.ok;
  file.status = true;
}

export function uploadFiles(files: UploadableFile[]) {
  return Promise.all(files.map((file) => uploadFile(file)));
}

export default function createUploader() {
  return {
    uploadFile: function (file: UploadableFile) {
      return uploadFile(file);
    },
    uploadFiles: function (files: UploadableFile[]) {
      return uploadFiles(files);
    },
  };
}
