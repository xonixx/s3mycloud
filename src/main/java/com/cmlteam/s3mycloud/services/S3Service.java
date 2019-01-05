package com.cmlteam.s3mycloud.services;

import com.cmlteam.s3mycloud.S3Props;
import com.cmlteam.util.Util;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Request;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Response;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.model.S3Object;

import javax.annotation.PostConstruct;
import java.io.File;
import java.net.URI;

@Service
@Slf4j
public class S3Service {
  private final S3Props s3Props;

  private S3Client s3Client;

  @Autowired
  public S3Service(S3Props s3Props) {
    this.s3Props = s3Props;
  }

  @PostConstruct
  public void postConstruct() {
    log.info("S3 Props: {}", s3Props);
    s3Client =
        S3Client.builder()
            .region(Region.of(s3Props.getRegion()))
            .endpointOverride(URI.create(s3Props.getEndpoint()))
            .credentialsProvider(
                () -> AwsBasicCredentials.create(s3Props.getAccessKey(), s3Props.getSecretKey()))
            .build();
    //    testUpload();
    testList();
  }

  private void testUpload() {
    log.info("Starting upload");
    long t0 = System.currentTimeMillis();
    String file = "/home/xonix/Downloads/haproxy_exporter-0.9.0.linux-amd64.tar.gz";
    String name = "haproxy_exporter-0.9.0.linux-amd64.tar.gz";
    s3Client.putObject(
        PutObjectRequest.builder().bucket(s3Props.getBucket()).key(name).build(),
        RequestBody.fromFile(new File(file)));
    log.info("Upload took: {}", Util.renderDurationFromStart(t0));
  }

  private void testList() {
    // maxKeys is set to 2 to demonstrate the use of
    // ListObjectsV2Result.getNextContinuationToken()
    int maxKeys = 50;
    ListObjectsV2Request req =
        ListObjectsV2Request.builder().bucket("nextcloud-xonix-1").maxKeys(maxKeys).build();
    ListObjectsV2Response resp;

    long totalSize = 0;

    do {
      resp = s3Client.listObjectsV2(req);

      for (S3Object objectSummary : resp.contents()) {
        Long size = objectSummary.size();
        totalSize += size;
        log.info(" - {} (size: {})", objectSummary.key(), Util.renderFileSize(size));
      }
      // If there are more than maxKeys keys in the bucket, get a continuation token
      // and list the next objects.
      String token = resp.nextContinuationToken();
      log.info("Next Continuation Token: {}", token);
      req =
          ListObjectsV2Request.builder()
              .bucket("nextcloud-xonix-1")
              .continuationToken(token)
              .maxKeys(maxKeys)
              .build();

    } while (resp.isTruncated());

    log.info("TOTAL: {}", Util.renderFileSize(totalSize));
  }
}
