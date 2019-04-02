FROM alpine:3.8
COPY ./build/bin/google-cloud-dns-healthcheck /bin/google-cloud-dns-healthcheck
ENTRYPOINT ["/bin/google-cloud-dns-healthcheck"]
CMD [ "run" ]
