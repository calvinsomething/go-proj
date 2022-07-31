FROM golang:1.18-alpine3.16
RUN adduser --system app
USER app
WORKDIR /app
CMD ["go", "run", "."]