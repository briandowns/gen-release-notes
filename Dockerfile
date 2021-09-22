FROM scratch

COPY bin/gen-release-notes /

CMD ["/gen-release-notes"]
