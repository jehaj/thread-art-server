FROM rust:slim AS builder
RUN apt-get update -qq && apt-get install -qq -y git > /dev/null
WORKDIR /usr/src
RUN git clone https://github.com/NikolajK-HTX/thread-art-rust.git
WORKDIR /usr/src/thread-art-rust
RUN cargo install --path .


FROM denoland/deno

ENV QUEUE_PATH="./thread-art-vol/queue"
ENV SAVE_PATH="./thread-art-vol/saves"
ENV EXECUTABLE_PATH="thread-art-rust"
ENV CONCURRENCY=2

WORKDIR /app
COPY --from=builder /usr/local/cargo/bin/thread-art-rust /usr/local/bin/thread-art-rust
RUN chmod -R 777 .

USER deno

COPY deps.ts .
RUN deno cache deps.ts

ADD . .
RUN deno cache main.ts

CMD [ "run", "--allow-env", "--allow-read", "--allow-write", "--allow-run", "main.ts" ]
