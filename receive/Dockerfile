FROM denoland/deno

RUN apt-get -qq update
RUN apt-get -qq -y install imagemagick > /dev/null

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

CMD [ "run", "--allow-net", "--allow-env", "--allow-read", "--allow-write", "--allow-run", "main.ts" ]
