import type { UploadableFile } from "./UploadableFile";
import { BACKEND_BASE_URL } from "@/urls";
import axios from "axios";

export async function uploadFile(file: UploadableFile) {
  // track status and upload file
  file.status = "loading";

  try {
    const response1: { data: { id: string; uploadUrl: string } } =
      await axios.post(BACKEND_BASE_URL + "/api/file/upload", {
        name: file.file.name,
        size: file.file.size,
        tags: ["TODO"],
      });

    const response2 = await axios.put(response1.data.uploadUrl, file.file, {
      headers: {
        "Content-Type": file.file.type,
      },
    });
  } catch (e) {
    console.error(e);
    file.status = false;
    return;
  }

  // change status to indicate the success of the upload request
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
