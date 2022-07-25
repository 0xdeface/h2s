FROM scratch
COPY h2s h2s
COPY ./internal/http/views/index.html /internal/http/views/index.html
COPY ./tpls ./tpls
EXPOSE 8090
CMD ["./h2s"]
