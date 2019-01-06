import axios from "axios";

class FilesApi {
  constructor(apiRoot) {
    this.apiRoot = apiRoot;
  }

  async ls(folder, next) {
    try {
      const resp = await axios.get(this.apiRoot + "/ls", {
        params: {
          folder,
          next,
          limit: 50
        }
      });
      return resp.data;
    } catch (e) {
      alert(e);
      throw e;
    }
  }
}

const filesApi = new FilesApi("http://localhost:8080"); // TODO

export default filesApi;
