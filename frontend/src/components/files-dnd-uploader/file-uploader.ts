import type { UploadableFile } from "./UploadableFile";
import { BACKEND_BASE_URL } from "@/urls";
import axios from "axios";

export async function uploadFile(file: UploadableFile, tags: string[]) {
  // track status and upload file
  file.status = "loading";

  try {
    const response1: { data: { id: string; uploadUrl: string } } =
      await axios.post(BACKEND_BASE_URL + "/api/file/upload", {
        name: file.file.name,
        size: file.file.size,
        tags,
      });

    try {
      const response2 = await axios.put(response1.data.uploadUrl, file.file, {
        headers: {
          "Content-Type": file.file.type,
        },
      });
      const response3 = await axios.post(
        BACKEND_BASE_URL + "/api/file/confirmUpload",
        {
          id: response1.data.id,
          success: true,
        }
      );
    } catch (e1: any) {
      const response3 = await axios.post(
        BACKEND_BASE_URL + "/api/file/confirmUpload",
        {
          id: response1.data.id,
          success: false,
          error: e1.message,
        }
      );
      throw e1;
    }
  } catch (e) {
    console.error(e);
    file.status = false;
    return;
  }

  // change status to indicate the success of the upload request
  file.status = true;
}

export function uploadFiles(files: UploadableFile[], tags: string[]) {
  return Promise.all(files.map((file) => uploadFile(file, tags)));
}

export default function createUploader() {
  return {
    uploadFiles: function (files: UploadableFile[], tags: string[]) {
      return uploadFiles(files, tags);
    },
  };
}
