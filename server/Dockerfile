FROM golang:1.18-alpine3.16
RUN addgroup --gid 1000 --system server
RUN adduser -u 1000 -G server --system server
WORKDIR /app
COPY . .
RUN go mod vendor && go mod verify
USER server
EXPOSE 8080
CMD ["go", "run", "."]