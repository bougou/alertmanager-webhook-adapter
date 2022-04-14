FROM golang:1.18

COPY / /src

RUN cd /src \
  && cd "cmd/alertmanager-webhook-adapter" \
  && go build -v -o /alertmanager-webhook-adapter \
  && true

ENTRYPOINT [ "/alertmanager-webhook-adapter" ]
CMD ["-v"]
