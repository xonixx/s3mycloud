package com.cmlteam.s3mycloud.services;

import com.amazonaws.HttpMethod;
import com.amazonaws.auth.AWSStaticCredentialsProvider;
import com.amazonaws.auth.BasicAWSCredentials;
import com.amazonaws.client.builder.AwsClientBuilder;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import com.amazonaws.services.s3.model.GeneratePresignedUrlRequest;
import com.cmlteam.s3mycloud.S3Props;
import com.cmlteam.s3mycloud.model.FileObj;
import com.cmlteam.s3mycloud.model.LsRequest;
import com.cmlteam.s3mycloud.model.LsResponse;
import com.cmlteam.util.Util;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.*;

import javax.annotation.PostConstruct;
import java.io.File;
import java.net.URI;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;

@Service
@Slf4j
public class S3Service {
  public static final int DEFAULT_LIMIT = 50;
  private final S3Props s3Props;

  private S3Client s3Client;
  private AmazonS3 s3ClientV1;

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

    s3ClientV1 =
        AmazonS3ClientBuilder.standard()
            .withEndpointConfiguration(
                new AwsClientBuilder.EndpointConfiguration(
                    s3Props.getEndpoint(), s3Props.getRegion()))
            .withCredentials(
                new AWSStaticCredentialsProvider(
                    new BasicAWSCredentials(s3Props.getAccessKey(), s3Props.getSecretKey())))
            .build();

    //    testUpload();
    //    testList();
    //    testListTopLevel();
    //    testPresigned();
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
    // maxKeys is set to demonstrate the use of
    // ListObjectsV2Result.getNextContinuationToken()
    int maxKeys = 50;
    ListObjectsV2Response resp;

    long totalSize = 0;

    String token = null;

    do {
      ListObjectsV2Request req =
          ListObjectsV2Request.builder()
              .bucket("nextcloud-xonix-1")
              .continuationToken(token)
              .maxKeys(maxKeys)
              .build();

      resp = s3Client.listObjectsV2(req);

      for (S3Object objectSummary : resp.contents()) {
        Long size = objectSummary.size();
        totalSize += size;
        log.info(" - {} (size: {})", objectSummary.key(), Util.renderFileSize(size));
      }
      // If there are more than maxKeys keys in the bucket, get a continuation token
      // and list the next objects.
      token = resp.nextContinuationToken();
      log.info("Next Continuation Token: {}", token);
    } while (resp.isTruncated());

    log.info("TOTAL: {}", Util.renderFileSize(totalSize));
  }

  private void testListTopLevel() {
    // maxKeys is set to demonstrate the use of
    // ListObjectsV2Result.getNextContinuationToken()
    int maxKeys = 50;
    ListObjectsV2Response resp;

    long totalSize = 0;

    String token = null;

    do {
      ListObjectsV2Request req =
          ListObjectsV2Request.builder()
              .bucket("nextcloud-xonix-1")
              .prefix("unindentified #2/")
              .delimiter("/")
              .continuationToken(token)
              .maxKeys(maxKeys)
              .build();

      resp = s3Client.listObjectsV2(req);

      for (CommonPrefix commonPrefix : resp.commonPrefixes()) {
        log.info(" - {}", commonPrefix.prefix());
      }
      for (S3Object objectSummary : resp.contents()) {
        Long size = objectSummary.size();
        totalSize += size;
        log.info(" - {} (size: {})", objectSummary.key(), Util.renderFileSize(size));
      }
      // If there are more than maxKeys keys in the bucket, get a continuation token
      // and list the next objects.
      token = resp.nextContinuationToken();
      log.info("Next Continuation Token: {}", token);
    } while (resp.isTruncated());

    log.info("TOTAL: {}", Util.renderFileSize(totalSize));
  }

  private void testPresigned() {
    // Set the presigned URL to expire after one hour.
    java.util.Date expiration = new java.util.Date();
    long expTimeMillis = expiration.getTime();
    expTimeMillis += 1000 * 60 * 60;
    expiration.setTime(expTimeMillis);

    // Generate the presigned URL.
    log.info("Generating pre-signed URL.");
    GeneratePresignedUrlRequest generatePresignedUrlRequest =
        new GeneratePresignedUrlRequest(s3Props.getBucket(), "test1.png")
            .withMethod(HttpMethod.PUT)
            .withExpiration(expiration);
    URL url = s3ClientV1.generatePresignedUrl(generatePresignedUrlRequest);

    log.info("Pre-Signed URL: {}", url.toString());
  }

  /**
   * @param lsRequest ls request params
   * @return list of files/folders in folder
   */
  public LsResponse ls(LsRequest lsRequest) {
    List<FileObj> items = new ArrayList<>();

    ListObjectsV2Response resp;

    String prefix = lsRequest.getFolder();
    if (prefix == null) prefix = "";
    if (!"".equals(prefix) && !prefix.endsWith("/")) prefix += "/";

    Integer limit = lsRequest.getLimit();
    if (limit == null || limit == 0) limit = DEFAULT_LIMIT;

    ListObjectsV2Request req =
        ListObjectsV2Request.builder()
            .bucket(s3Props.getBucket())
            .prefix(prefix)
            .delimiter("/")
            .continuationToken(lsRequest.getNext())
            .maxKeys(limit)
            .build();

    resp = s3Client.listObjectsV2(req);

    for (CommonPrefix commonPrefix : resp.commonPrefixes()) {
      //      log.info(" - {}", commonPrefix.prefix());
      String key = commonPrefix.prefix();
      String[] parts = key.split("/");
      items.add(FileObj.builder().name(parts[parts.length - 1]).isFolder(true).size("").build());
    }
    for (S3Object objectSummary : resp.contents()) {
      Long size = objectSummary.size();
      String key = objectSummary.key();
      String[] parts = key.split("/");
      String name = parts[parts.length - 1];
      if (("/" + prefix).endsWith("/" + name + "/") && size == 0L) continue; // same folder that we list
      items.add(
          FileObj.builder().name(name).isFolder(false).size(Util.renderFileSize(size)).build());
      //      log.info(" - {} (size: {})", objectSummary.key(), Util.renderFileSize(size));
    }
    String token = resp.nextContinuationToken();
    return new LsResponse(items, token);
  }
}
