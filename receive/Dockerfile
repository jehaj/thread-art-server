FROM denoland/deno:alpine

RUN apk -q update
RUN apk -q add imagemagick

ENV PORT=8001
ENV QUEUE_PATH="/thread-art-vol/queue"
ENV SAVE_PATH="/thread-art-vol/saves"
ENV QUEUE_LIMIT=10

EXPOSE $PORT
WORKDIR /app
RUN chmod -R 777 .

USER deno

COPY deps.ts .
RUN deno cache deps.ts

ADD . .
RUN deno cache main.ts

CMD [ "run", "-A", "--unstable", "main.ts" ]
