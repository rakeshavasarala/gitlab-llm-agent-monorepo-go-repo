FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/security-labeler

FROM alpine:3.18
COPY --from=build /bin/security-labeler /bin/security-labeler
ENTRYPOINT ["/bin/security-labeler"]
