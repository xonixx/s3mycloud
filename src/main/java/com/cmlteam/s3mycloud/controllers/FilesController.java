package com.cmlteam.s3mycloud.controllers;

import com.cmlteam.s3mycloud.model.LsRequest;
import com.cmlteam.s3mycloud.model.LsResponse;
import com.cmlteam.s3mycloud.services.S3Service;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.validation.Valid;

@RestController
@CrossOrigin(origins = "http://localhost:3000")
public class FilesController {
  private final S3Service s3Service;

  @Autowired
  public FilesController(S3Service s3Service) {
    this.s3Service = s3Service;
  }

  @GetMapping(value = "ls")
  public LsResponse ls(@Valid LsRequest lsRequest) {
    return s3Service.ls(lsRequest);
  }
}
