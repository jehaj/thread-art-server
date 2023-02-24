async function main() {
    console.log("hello world!");
    const data = await fetch("/example.txt");
    const result = (await data.text()).split('\n')[1].split(',');

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

main();
