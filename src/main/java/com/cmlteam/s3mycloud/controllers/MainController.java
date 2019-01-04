package com.cmlteam.s3mycloud.controllers;

import com.cmlteam.s3mycloud.model.ServerStatus;
import com.cmlteam.s3mycloud.services.SampleService;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.RestTemplate;

import java.util.HashMap;
import java.util.Map;

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

  @RequestMapping(value = "/testws", method = RequestMethod.GET)
  public String testws() {
    RestTemplate template = new RestTemplate();
    ServerStatus status =
        template.getForObject("https://l2c1x1.com/services/misc/server-stats", ServerStatus.class);
    return "" + status.totalAccounts;
  }

  @RequestMapping(value = "/testdb", method = RequestMethod.GET)
  public String testdb() {
    return sampleService.getDbVersion();
  }

  @PostMapping(value = "/testpost")
  public Map testPost(@RequestBody Map payload) {
    log.info("Test POST: {}", payload);
    HashMap res = new HashMap();
    res.put("success", true);
    return res;
  }
}
