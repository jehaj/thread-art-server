# thread-art-server
This is to be used in combination with `thread-art-rust` or another software 
that makes (as is seen in art by Petros Vrellis) an image of threads by
a user uploaded image. It contains to projects that work together as one.
The `receive` part receives an image by the user and adds it to the queue.
Then the `<to-be-named>` part runs the algorithm on one (or more) image at a 
time from the queue.

It is a rework of the work done at `thread-art-archive`. A lot of work was 
done and can be found there - including the algorithm in different languages.

## Specifics
ImageMagick is used to convert a user uploaded image into the format that the
software accepts. The default case is 400x400 grayscale:
```
convert <input>.png -resize 400x400^ -gravity center -extent 400x400 -colorspace Gray <output>.png
```

The goal is to use `Docker` (`podman`) to run these two parts as two 
containers that communicate by using the file system as a queue. A 
`docker volume` is used as the file system, where both have access.

There are different ways of going about using `thread-art-rust`:
- Using `git clone` and compiling it in the `Dockerfile`.
- Having compiled it before and downloading it with `ADD`.

I am leaning towards the first option. After some research it can be done with
[multi-stage builds in docker](https://docs.docker.com/build/building/multi-stage/).

## Commands
```
$ deno run --allow-net --allow-write --allow-read --allow-env --allow-run main.ts
```

```
$ podman-compose up
```

```
$ docker build -t receive . && docker run -p 8001:8001 receive
```

```
$ podman kube play kube.yaml
$ docker build -t receive .
$ podman kube play --replace kube.yaml
```
