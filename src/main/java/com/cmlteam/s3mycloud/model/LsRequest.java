package com.cmlteam.s3mycloud.model;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class LsRequest {
  private String folder;
  private String next;
  private Integer limit;
}
