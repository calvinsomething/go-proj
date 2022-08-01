FROM golang:1.18-alpine3.16
RUN adduser --system app
USER app
WORKDIR /app
EXPOSE 8080
RUN ls
CMD ["go", "run", "."]