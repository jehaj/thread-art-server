let params = (new URL(document.location)).searchParams;
let id = params.get("id");

const idHolder = document.getElementById("id-holder");
idHolder.innerText = id;

let timeout = 1000;

const apiURL = window.location.origin + "/api";

function getImage() {
  fetch(`${apiURL}/${id}`).then((res) => {
    const resultText = document.getElementById("result-text");
    resultText.removeAttribute("hidden");
    if (res.status == 204) {
      resultText.innerText = "The algorithm can not run correctly on this image. Try another.";
      image.setAttribute("src", window.origin + "/error.png");
      image.removeAttribute("hidden");
      return;
    }
    console.log(res);
    if (res.status == 200) {
      resultText.innerText = "Success! See your image:";
      const image = document.getElementById("image");
      res.formData().then((data) => {
        const imageUrl = URL.createObjectURL(data.get("image"));
        image.setAttribute("src", imageUrl);
        image.removeAttribute("hidden");

        // show text carousel
        const textData = data.get("text");
        updateCarousel(textData);
      });
      deleteBtn.removeAttribute("hidden");
    } else {
      resultText.innerText = "Failed. Trying again :)";
      setTimeout(getImage, timeout);
      timeout = 2 * timeout;
    }
  });
}

async function updateCarousel(data) {
  const result = data.split('\n')[1].split(',');

  const carousel = document.getElementById("carousel");
  for (let i = 0; i < result.length - 1; i++) {
    const step = document.createElement("div");
    const point = result[i];
    step.id = "step-" + (i + 1);
    step.innerText = point;
    const stepnumber = document.createElement("p");
    stepnumber.innerText = (i + 1);
    step.appendChild(stepnumber);
    carousel.appendChild(step);
  }
}

getImage();
const deleteInfo = document.getElementById("delete-info");
const deleteBtn = document.getElementById("delete-btn");
deleteBtn.addEventListener("click", deleteImage);

function deleteImage() {
  fetch(`${apiURL}/${id}`, {
    method: "DELETE"
  }).then((res) => {
    if (res.status == 200) {
      window.location.assign(window.location.origin);
    } else {
      deleteInfo.innerHTML = "Error. The image was not found. Is it already deleted?";
    }
  })
}
