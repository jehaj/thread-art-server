# thread-art-server
This is to be used in combination with `thread-art-rust` or another software 
that makes (as is seen in art by Petros Vrellis) an image of threads by
a user uploaded image. It contains multiple projects that work together as one.

It is a rework of the work done at `thread-art-archive`. A lot of work was 
done and can be found there - including the algorithm in different languages.

Try visiting and using it at [art.jehaj.dk](https://art.jehaj.dk).

## view
`view` is the frontend part of `thread-art-server`. It uses
- vuejs
- bulma
`typescript` is also used.

## backbone
Is the API server. It is written in `go`. It uses
- gorm
- chi
- sqlite
- [thread-art](https://github.com/jehaj/thread-art-rust) (the algorithm)
`bruno` is used to check the API endpoints.

## Development tools
I recommend using scoop to setup various tools. Visit 
[scoop.sh](https://scoop.sh) and setup scoop. Check with 
```
scoop checkup
```
that everything is working as expected. Then use
```
scoop install nodejs go bruno pnpm
```
to get the tools. The last one is optional. You might need to add a bucket for
`bruno`. This is done with
```
scoop bucket add extras
```
Visit [scoop.sh](https://scoop.sh) if trouble arises.


I use JetBrains editors, but again that is optional. WebStorm is used for 
`lens` and GoLand for `backbone`.

Further details can be found in the respective folders `README.md`.
