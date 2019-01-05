import React, { Component } from "react";
import filesApi from "../api/FilesApi";

class Files extends Component {
  state = {
    selectedFolder: "",
    items: [],
    next: null
  };

  constructor(params) {
    super(params);
  }

  async componentDidMount() {
    await this.load();
  }

  async load() {
    const { items, next } = await filesApi.ls(this.state.selectedFolder, null);
    this.setState({ items, next });
  }

  async loadNext() {}

  async loadFolder(name) {
    this.setState({ selectedFolder: `${this.state.selectedFolder}/${name}`, next: null });
    await this.load();
  }

  render() {
    return (
      <div>
        {this.state.items.map(f => (
          <div style={{ border: "dotted 1px #ddd" }}>
            {f.folder ? (
              <a href={""} onClick={() => this.loadFolder(f.name)}>
                {f.name}
              </a>
            ) : (
              f.name
            )}
            {f.size ? `(${f.size})` : ""}
          </div>
        ))}
        {this.state.next && (
          <a href={""} onClick={this.loadNext}>
            more...
          </a>
        )}
      </div>
    );
  }
}

export default Files;
