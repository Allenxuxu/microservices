FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY user-api /user-api
ENV CONSUL_ADDRESS=127.0.0.1:8500 \
    MICRO_REGISTRY=consul \
    MICRO_REGISTRY_ADDRESS=127.0.0.1:8500
ENTRYPOINT /user-api
LABEL Name=user-api Version=0.0.1
