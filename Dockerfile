FROM golang:1.18

COPY / /src

RUN cd /src \
  && make build \
  && mv /src/_output/alertmanager-webhook-adapter /alertmanager-webhook-adapter \
  && rm -rf /src \
  && true

ENTRYPOINT [ "/alertmanager-webhook-adapter" ]
CMD ["-v"]
