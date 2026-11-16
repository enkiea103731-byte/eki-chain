FROM node:20-alpine
WORKDIR /app
COPY package.json ./
RUN npm install
COPY site ./site
COPY contracts.json eki.json ./
EXPOSE 10000
CMD ["npx","serve","site","-l","10000"]
