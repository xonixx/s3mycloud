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

  load = async () => {
    const { items: moreItems, next } = await filesApi.ls(
      this.state.selectedFolder,
      this.state.next
    );
    this.setState({ items: [...this.state.items, ...moreItems], next });
  };

  loadFolder = name => {
    this.setState(
      { items: [], selectedFolder: `${this.state.selectedFolder}${name}`, next: null },
      this.load
    );
  };

  render() {
    return (
      <div style={{ padding: 10 }}>
        {this.state.items.map(f => (
          <div key={f.name} style={{ border: "dotted 1px #ddd" }}>
            {f.folder ? (
              <a
                href={"#"}
                onClick={e => {
                  e.preventDefault();
                  this.loadFolder(f.name);
                }}
              >
                {f.name}
              </a>
            ) : (
              f.name
            )}
            {f.size ? `(${f.size})` : ""}
          </div>
        ))}
        {this.state.next && (
          <a
            href={"#"}
            onClick={e => {
              e.preventDefault();
              this.load();
            }}
          >
            more...
          </a>
        )}
      </div>
    );
  }
}

export default Files;
