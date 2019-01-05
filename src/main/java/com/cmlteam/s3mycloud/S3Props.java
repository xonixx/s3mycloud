package com.cmlteam.s3mycloud;

import lombok.Getter;
import lombok.Setter;
import lombok.ToString;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "s3")
@Getter
@Setter
@ToString(exclude = {"accessKey", "secretKey"})
public class S3Props {
  private String bucket;
  private String accessKey;
  private String secretKey;
  private String region;
  private String endpoint;
}
