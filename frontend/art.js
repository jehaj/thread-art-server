let params = (new URL(document.location)).searchParams;
let id = params.get("id");

const idHolder = document.getElementById("id-holder");
idHolder.innerText = id;

let timeout = 1000;

function getImage() {
    fetch(`http://localhost:8001/${id}`).then((res) => {
        const successHolder = document.getElementById("success-holder");
        successHolder.removeAttribute("hidden");
        console.log(res);
        if (res.status == 200) {
            successHolder.innerText = "Success! See your image:";
            const image = document.getElementById("image");
            res.formData().then((data) => {
                const imageUrl = URL.createObjectURL(data.get("image"));
                image.removeAttribute("hidden");
                image.setAttribute("src", imageUrl);
            })
        } else {
            successHolder.innerText = "Failed. Trying again :)";
            setTimeout(getImage, timeout);
            timeout = 2 * timeout;
        }
    });
}

getImage();
