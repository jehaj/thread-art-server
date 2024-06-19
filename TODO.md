# TODO
What is missing? Some things...

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
- [x] send back original image
- [x] send back points.txt (as json array)
- [x] send back image IDs with user ID.
- [x] change status to finished true when done

## view
What you see and can click on to interact with the project. It uses the API. It 
uses bulma, vuejs and typescript.
- [x] homepage
  - [x] user id (fill automatically if in localStorage)
  - [x] upload field
  - [x] description
- [x] personal page
  - [x] overview with thumbnails
  - [x] perhaps animate on mouseover
- [ ] image page
  - [x] thumbnail
  - [x] thread-art-animation (with canvas)
  - [ ] download .txt file / image
    - [ ] points.txt has the points
  - [ ] download image
  - [ ] scrollwheel with points

Der skal være en eksempel bruger (demo), så man ikke selv behøver uploade et
billede for at se hvad det går ud på. Jeg tænker, at denne bruger skal laves
første gang go-serveren kører starter (den tjekker så om demo brugeren 
allerede eksisterer, hvis ikke laver det den). F.eks. har den bruger ID 
"`1-red-dragon`" og har tre billeder. Det skal selvfølgelig være kodet på en
måde, så man kan undgå at tilføje denne demo-bruger. Det skal *lens* også
understøtte, at demo-brugeren kan være der og ikke være der.

DEMO_USER_ID != "", så findes der en demo-bruger. Eller også skal det være et 
argument til kørsel-filen. At det er en miljøvariabel gør at f.eks. det kun skal
skrives et sted og begge kan så bruge det (i f.eks. Dockerfile). I WebStorm og
GoLand skal det skrives begge steder.
