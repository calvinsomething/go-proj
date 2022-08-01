FROM node:18-alpine3.15
RUN adduser --system app
USER app
WORKDIR /app
COPY package*.json .
RUN npm i
EXPOSE 3000
CMD ["npm", "start"]