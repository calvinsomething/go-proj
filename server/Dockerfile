FROM golang:1.18-alpine3.16
RUN addgroup --gid 1000 --system server
RUN adduser -u 1000 -G server --system server
USER server
WORKDIR /app
EXPOSE 8080
CMD ["go", "run", "."]