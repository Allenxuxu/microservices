FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY user-srv /user-srv
ENV CONSUL_ADDRESS=127.0.0.1:8500 \
    MICRO_REGISTRY=consul \
    MICRO_REGISTRY_ADDRESS=127.0.0.1:8500
ENTRYPOINT /user-srv
LABEL Name=user-srv Version=0.0.1
