package com.cmlteam.s3mycloud.controllers;

import com.cmlteam.s3mycloud.services.SampleService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
@Slf4j
public class MainController {

  private final SampleService sampleService;

  @Autowired
  public MainController(SampleService sampleService) {
    this.sampleService = sampleService;
  }

  @RequestMapping(value = "/", method = RequestMethod.GET)
  public String test() {
    return "Hello CML!";
  }
}
