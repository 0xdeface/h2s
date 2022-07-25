FROM scratch
COPY h2s h2s
EXPOSE 8090
RUN mkdir -p "/internal/http/views/"
COPY ./internal/http/views/index.html /internal/http/views/index.html
ADD ./tpls ./tpls
CMD ["./h2s"]
