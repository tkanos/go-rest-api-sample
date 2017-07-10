FROM scratch
COPY app /
ENTRYPOINT ["./app"]
EXPOSE 8001