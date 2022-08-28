FROM node:16-alpine3.15
RUN deluser --remove-home node
RUN adduser --uid 1000 -D client
WORKDIR /app
USER client
EXPOSE 3000
CMD ["npm", "start"]
