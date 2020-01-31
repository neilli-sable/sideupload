# sideupload

## About
- docker-compose とかで、データボリュームをS3にバックアップするために使う

## Example

```
services:
  sideupload:
    image: neilli/sideupload
    command: cron
    environment:
      SIDEUPLOAD_LISTENPORT: 8888
      SIDEUPLOAD_TARGETDIR: "/backup"
      SIDEUPLOAD_CRONWITHSECOND: "0 * * * * *"
      SIDEUPLOAD_BACKUPSTORAGE_CUSTOMENDPOINT: "http://localhost:9000"
      SIDEUPLOAD_BACKUPSTORAGE_REGION: "us-west-2"
      SIDEUPLOAD_BACKUPSTORAGE_PREFIX: ""
      SIDEUPLOAD_BACKUPSTORAGE_BUCKETNAME: "example-backet"
      SIDEUPLOAD_BACKUPSTORAGE_ACCESSKEY: minio
      SIDEUPLOAD_BACKUPSTORAGE_SECRETKEY: minio_good_secretkey
    networks:
      - sideupload_net
    volumes:
      - "../samples:/backup/samples"
      - "../samples:/backup/samples1"
```