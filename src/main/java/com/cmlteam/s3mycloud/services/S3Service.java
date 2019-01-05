package com.cmlteam.s3mycloud.services;

import com.cmlteam.s3mycloud.S3Props;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;

@Service
@Slf4j
public class S3Service {
  private final S3Props s3Props;

  @Autowired
  public S3Service(S3Props s3Props) {
    this.s3Props = s3Props;
  }

  @PostConstruct
  public void postConstruct() {
    log.info("S3 Props: {}", s3Props);
  }
}
