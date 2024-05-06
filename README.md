# thread-art-server
This is to be used in combination with `thread-art-rust` or another software 
that makes (as is seen in art by Petros Vrellis) an image of threads by
a user uploaded image. It contains multiple projects that work together as one.

It is a rework of the work done at `thread-art-archive`. A lot of work was 
done and can be found there - including the algorithm in different languages.

## view
`view` is the frontend part of `thread-art-server`. It uses
- vuejs
- bulma

## backbone
Is the API server. It uses
- go
- gorm
- mux
