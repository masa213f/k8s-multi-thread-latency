FROM scratch

COPY burn /burn

USER 10000:10000

ENTRYPOINT ["/burn"]
