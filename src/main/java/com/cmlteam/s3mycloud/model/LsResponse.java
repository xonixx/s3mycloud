package com.cmlteam.s3mycloud.model;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.List;

@Getter
@AllArgsConstructor
public class LsResponse {
  private final List<FileObj> items;
  private final String next;
}
