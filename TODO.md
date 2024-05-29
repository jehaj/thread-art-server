# TODO
What is missing? Many things...

## backbone
The backend / API. Written in go.
- [x] upload image (save image)
- [x] convert image to grayscale
- [ ] add image to queue
- [ ] create user ID (number + color + animal)
  - [ ] if ID already exists then INTERNAL SERVER ERROR
- [ ] link ID and image in database (image gets UUID, but user gets funny ID)
- [ ] send back user id (which should be stored automatically for the user - can be ID given by user if exists)
- [ ] queue to worker pool
- [ ] worker uses thread-art-rust
- [ ] send back image
- [ ] send back points.txt (as json array?)

## view
What you see and can click on to interact with the project. It uses the API. It uses bulma, vuejs and typescript.
- [ ] homepage
  - [ ] user id (fill automatically if in cookie)
  - [ ] upload field
- [ ] personal page
  - [ ] overview with thumbnails
- [ ] image page
  - [ ] thumbnail
  - [ ] thread-art-animation (with canvas)
    - [ ] points.txt has the points
  - [ ] download image
  - [ ] scrollwheel with points