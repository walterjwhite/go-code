FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /src/cmd/example/website
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/website ./...

FROM nginx:alpine

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

COPY --from=builder /bin/website /bin/website
RUN chown -R appuser:appgroup /bin/website

COPY cmd/example/website/.build/nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /usr/share/nginx/html

COPY cmd/example/website/target/classes/* /usr/share/nginx/html

RUN chown -R nginx:nginx /var/cache/nginx /var/log/nginx /etc/nginx /usr/share/nginx/html && \
    touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

COPY cmd/example/website/.build/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

USER nginx

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8180/ || exit 1

EXPOSE 8180

CMD ["/docker-entrypoint.sh"]
