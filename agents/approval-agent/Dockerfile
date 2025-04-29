FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/approval-agent

FROM alpine:3.18
COPY --from=build /bin/approval-agent /bin/approval-agent
ENTRYPOINT ["/bin/approval-agent"]
