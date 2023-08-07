FROM --platform=linux/amd64 node:16-alpine AS build
# 이미지를 node:10으로 당겨왔다가 1시간동안 뻘짓함
# 도커 이미지 빌드 중 npm run build에서 막히길래 한참 찾았음

RUN mkdir app

WORKDIR /app

COPY ./package.json ./

RUN npm install
# package.json을 해야 코드의 일부를 수정할 때 매번 종속성까지 다시 다운받지 않을 수 있음

COPY . .

RUN npm run build
# axios 등 다른 패키지를 사용하고 도커 이미지를 만드려면 package.json에 해당 패키지를 작성해줘야 함

FROM --platform=linux/amd64 nginx

EXPOSE 3000

COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf

COPY --from=build /app/build /usr/share/nginx/html 