FROM scratch

COPY ./deploy /usr/src/app/

# Use an unprivileged user.
USER appuser
WORKDIR /usr/src/app

# Run binary.
ENTRYPOINT ["/usr/src/app/deploy"]