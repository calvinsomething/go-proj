FROM node:16-alpine3.15
# Rename node user
RUN sed -i 's/node:/client:/' /etc/passwd /etc/group /etc/shadow
WORKDIR /app
USER client
EXPOSE 3000
CMD ["npm", "start"]
