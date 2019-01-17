#
# Dockerfile - Git Archivist (git-archivist).
#

FROM alpine:latest
LABEL vendor="sostheim"

RUN apk update && \
    apk add git 

COPY build/linux_amd64/git-archivist /

ENTRYPOINT ["/git-archivist"]
