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

  loadFolder = name => this.loadFolderFull(`${this.state.selectedFolder}${name}/`);

  loadFolderFull = fullPath => {
    this.setState({ items: [], selectedFolder: fullPath, next: null }, this.load);
  };

  render() {
    const selectedFolder = this.state.selectedFolder;
    let selectedFolderParts = selectedFolder.split("/");
    selectedFolderParts = selectedFolderParts.slice(0, selectedFolderParts.length - 1);
    selectedFolderParts = ["", ...selectedFolderParts];
    return (
      <div style={{ padding: 10 }}>
        {selectedFolder && (
          <div style={{ marginBottom: 10 }}>
            {selectedFolderParts.map((part, i) => {
              let fullPath = selectedFolderParts.slice(1, i + 1).join("/") + "/";
              if (fullPath === "/") fullPath = "";
              return (
                <span key={fullPath}>
                  {i === selectedFolderParts.length - 1 ? (
                    part
                  ) : (
                    <span>
                      <a
                        href={""}
                        title={fullPath}
                        onClick={e => {
                          e.preventDefault();
                          this.loadFolderFull(fullPath);
                        }}
                      >
                        {part || "home"}
                      </a>
                      {" / "}
                    </span>
                  )}
                </span>
              );
            })}
          </div>
        )}
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
            {f.size ? ` (${f.size})` : ""}
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
