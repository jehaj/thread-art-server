# TODO
What is missing? Many things...

## backbone
The backend / API. Written in go.
- [x] upload image (save image)
- [x] convert image to grayscale
- [x] add image to queue
- [x] create user ID (number + color + animal)
  - [x] if ID already exists then INTERNAL SERVER ERROR
- [x] link ID and image in database (image gets UUID, but user gets funny ID)
- [x] send back user id (which should be stored automatically for the user - can be ID given by user if exists)
- [x] queue to worker pool
- [x] worker uses thread-art-rust
- [x] send back image
- [x] send back points.txt (as json array)
- [x] send back image IDs with user ID.

## view
What you see and can click on to interact with the project. It uses the API. It uses bulma, vuejs and typescript.
- [ ] homepage
  - [ ] user id (fill automatically if in cookie)
  - [ ] upload field
  - [ ] description
- [ ] personal page
  - [ ] overview with thumbnails
  - [ ] perhaps animate on mouseover
  - [ ] download .txt file / image
- [ ] image page
  - [ ] thumbnail
  - [ ] thread-art-animation (with canvas)
    - [ ] points.txt has the points
  - [ ] download image
  - [ ] scrollwheel with points
