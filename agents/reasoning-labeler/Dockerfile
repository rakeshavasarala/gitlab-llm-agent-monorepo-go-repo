FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/reasoning-labeler

FROM alpine:3.18
COPY --from=build /bin/reasoning-labeler /bin/reasoning-labeler
ENTRYPOINT ["/bin/reasoning-labeler"]
