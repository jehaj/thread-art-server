FROM golang:bookworm AS builder

WORKDIR /usr/src/app

COPY backbone/go.mod backbone/go.sum ./
RUN go mod download && go mod verify 

COPY backbone/ ./
RUN go build .


FROM node:bookworm AS frontend

WORKDIR /usr/src/app

COPY lens/package.json ./
RUN npm install

COPY lens/ ./
ENV API_URL=http://localhost:8080
ENV VITE_DEMO_USER_ID=1-yellow-dragon
RUN npm run build


FROM debian:bookworm-slim

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/backbone ./

COPY --from=builder /usr/src/app/color-list.txt ./
COPY --from=builder /usr/src/app/animal-list.txt ./

COPY --from=builder /usr/src/app/demo1.png ./
COPY --from=builder /usr/src/app/demo2.png ./
COPY --from=builder /usr/src/app/demo3.png ./
COPY --from=frontend /usr/src/app/dist ./dist 
COPY backbone/thread-art-rust ./

ENV DEMO_USER_ID=1-yellow-dragon

CMD ./backbone

