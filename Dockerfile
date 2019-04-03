FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY ./build/bin/google-cloud-dns-healthcheck /bin/google-cloud-dns-healthcheck
ENTRYPOINT ["/bin/google-cloud-dns-healthcheck"]
CMD [ "run" ]
