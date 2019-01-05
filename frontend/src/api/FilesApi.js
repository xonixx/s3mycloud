import axios from "axios";

class FilesApi {
  constructor(apiRoot) {
    this.apiRoot = apiRoot;
  }

  async ls(folder, next) {
    return (await axios.get(this.apiRoot + "/ls", {
      params: {
        folder,
        next,
        limit: 50
      }
    })).data;
  }
}

const filesApi = new FilesApi("localhost:8080"); // TODO

export default filesApi;
