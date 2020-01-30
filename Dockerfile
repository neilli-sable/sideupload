FROM neilli/golang-with-webpack AS build-env
ADD . /workdir
WORKDIR /workdir
RUN make build && make compress

FROM alpine

RUN apk add --no-cache ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN mkdir -p /var/www/app
COPY --from=build-env /workdir/build/linux-amd64/bin/sideupload /var/www/app/sideupload
COPY --from=build-env /workdir/sideupload.toml /var/www/app/sideupload.toml

WORKDIR /var/www/app/
ENTRYPOINT [ "./sideupload" ]
