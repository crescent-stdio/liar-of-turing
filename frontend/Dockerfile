FROM node:latest as base

WORKDIR /usr/frontend

COPY package*.json ./

RUN npm install -g pnpm

RUN pnpm install

COPY . .

ENV NODE_ENV=production

RUN pnpm run build

EXPOSE 3000

CMD ["pnpm", "start"]