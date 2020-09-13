FROM alpine:latest

# set labels for metadata
LABEL name="cosi-defaultbc" \
  description="A Kubernetes mutating webhook server that implements default bucketclass" \
  summary="A Kubernetes mutating webhook server that implements default bucketclass"

# install sidecar-injector binary
COPY bin/cosi-webhook /cosi-webhook

# set entrypoint
ENTRYPOINT ["/cosi-webhook"]

