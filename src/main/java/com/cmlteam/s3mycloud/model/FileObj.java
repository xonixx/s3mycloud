package com.cmlteam.s3mycloud.model;

import lombok.Builder;
import lombok.Getter;

@Getter
@Builder
public class FileObj {
  private boolean isFolder;
  private String name;
  /** Rendered size str */
  private String size;
}
