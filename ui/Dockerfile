FROM node:18

WORKDIR /app   
COPY package*.json ./
RUN npm ci  
 
COPY . .
RUN npm run build  
ENV TZ=UTC

EXPOSE 80
CMD ["npm", "start"] 