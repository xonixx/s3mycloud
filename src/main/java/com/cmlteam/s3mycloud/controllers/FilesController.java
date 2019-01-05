package com.cmlteam.s3mycloud.controllers;

import com.cmlteam.s3mycloud.services.S3Service;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class FilesController {
    private final S3Service s3Service;

    @Autowired
    public FilesController(S3Service s3Service) {
        this.s3Service = s3Service;
    }


}
