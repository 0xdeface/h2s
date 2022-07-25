FROM scratch
COPY h2s h2s
EXPOSE 8090
CMD ["./h2s"]
