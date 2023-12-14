FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-verkada"]
COPY baton-verkada /