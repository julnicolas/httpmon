FROM golang:1.21-alpine AS build

WORKDIR /app
COPY . /app

RUN go build -o ./bin/httpmon main.go

FROM golang:1.21-alpine

WORKDIR /app

RUN addgroup httpmon_g && \
adduser httpmon -D -G httpmon_g

COPY --from=build --chown=httpmon:httpmon_g /app/bin/httpmon /app/
COPY --from=build --chown=httpmon:httpmon_g /app/sample_csv.txt /app/

USER httpmon
ENTRYPOINT ["./httpmon"]
CMD ["--file", "./sample_csv.txt"]
