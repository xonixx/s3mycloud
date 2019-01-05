#!/usr/bin/env bash

#curl -s -XPUT 'https://s3mycloud-1.s3.us-west-1.wasabisys.com/test1.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20190105T154814Z&X-Amz-SignedHeaders=host&X-Amz-Expires=3599&X-Amz-Credential=UC6LXH4ABQQ096627RM7%2F20190105%2Fus-west-1%2Fs3%2Faws4_request&X-Amz-Signature=846425f58562cd3ff595b6cd986659426baecc82a50253132f5f9635f3d4de2f' \
curl -s -XPUT 'https://s3mycloud-1.s3.us-west-1.wasabisys.com/test1.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20190105T154902Z&X-Amz-SignedHeaders=host&X-Amz-Expires=3599&X-Amz-Credential=UC6LXH4ABQQ096627RM7%2F20190105%2Fus-west-1%2Fs3%2Faws4_request&X-Amz-Signature=ccba41ff677663211e9bd01f305e5fcbad5ee269deef09a73897a025defd5e2e' \
    -H 'Content-Type: image/png' \
    --data-binary "@/home/xonix/Desktop/test1.png"

