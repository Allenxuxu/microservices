FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY hello-srv /hello-srv
ENV MICRO_REGISTRY=consul \
    MICRO_REGISTRY_ADDRESS=127.0.0.1:8500
ENTRYPOINT /hello-srv
LABEL Name=hello-srv Version=0.0.1