package com.cmlteam.s3mycloud;

import lombok.Getter;
import lombok.Setter;
import lombok.ToString;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;
import org.springframework.validation.annotation.Validated;

import javax.validation.constraints.NotBlank;

@Component
@ConfigurationProperties(prefix = "s3")
@Validated
@Getter
@Setter
@ToString(exclude = {"accessKey", "secretKey"})
public class S3Props {
  @NotBlank private String bucket;
  @NotBlank private String accessKey;
  @NotBlank private String secretKey;
  @NotBlank private String region;
  @NotBlank private String endpoint;
}
