import { join, loadSync, Queue, serve } from "./deps.ts";

loadSync();

function RandomID() {
  const letters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const clen = letters.length;
  const r = crypto.getRandomValues(new Uint8Array(12));
  const b = [];
  for (let i = 0; i < 12; i++) {
    const c = r[i];
    b.push(letters[c % clen]);
  }
  return b.join("");
}

const QUEUE_LIMIT = parseInt(Deno.env.get("QUEUE_LIMIT") || "10");
const SAVE_PATH = Deno.env.get("SAVE_PATH") || "./saves";
const PORT = parseInt(Deno.env.get("PORT") || "8001");
const jobQueue = new Queue("job", {
  connection: {
    host: "localhost",
  },
  defaultJobOptions: {
    removeOnFail: true,
    attempts: 1,
    removeOnComplete: true,
  },
});

console.log("Running with settings:");
console.log("QUEUE_LIMIT =", QUEUE_LIMIT);
console.log("SAVE_PATH =", SAVE_PATH);
console.log("PORT =", PORT);
console.log("Working at:", Deno.cwd());

await Deno.mkdir(SAVE_PATH, { recursive: true });

serve(handler, { port: PORT });

async function handler(req: Request): Promise<Response> {
  if (req.method == "GET") {
    try {
      const id = req.url.substring(req.url.lastIndexOf("/") + 1);
      const f = await Deno.readTextFile(join(SAVE_PATH, id, "RESULT.txt"));
      const imageFile = await Deno.readFile(join(SAVE_PATH, id, "RESULT.png"));
      const data = new FormData();
      data.append("image", new Blob([imageFile], { type: "image/png" }));
      data.append("text", f);
      return new Response(data, { status: 200 });
    } catch (_) {
      try {
        const id = req.url.substring(req.url.lastIndexOf("/") + 1);
        (await Deno.open(join(SAVE_PATH, id, "ERROR"))).close();
        return new Response(null, { status: 204 });
      } catch (_) {
        // do nothing
      }
      return new Response(
        "The image at entered ID (if it exists) is not done yet.",
        { status: 404 },
      );
    }
  } else if (req.method == "POST") {
    const reqLength = req.headers.get("Content-Length");
    if (reqLength == null) {
      return new Response("Not allowed", { status: 400 });
    }
    if (parseInt(reqLength) > 200000) {
      return new Response(
        "Image is too large! Try uploading a smaller image.",
        { status: 500 },
      );
    }
    if (jobQueue.count() > QUEUE_LIMIT) {
      return new Response("Try again later! Queue is full.", { status: 503 });
    }
    try {
      const f = await req.formData();
      const imageEntry = f.get("image");
      if (imageEntry == null || typeof imageEntry != "object") {
        throw new Error("Image not contained in form data.");
      }
      const filename = RandomID();
      const image = new Uint8Array(await imageEntry.arrayBuffer());
      await Deno.mkdir(join(SAVE_PATH, filename), { recursive: true });
      await Deno.writeFile(join(SAVE_PATH, filename, "INPUT.png"), image, {
        createNew: true,
      });
      const cmd = [
        "convert",
        join(SAVE_PATH, filename, "INPUT.png"),
        "-resize",
        "400x400^",
        "-gravity",
        "center",
        "-extent",
        "400x400",
        "-colorspace",
        "Gray",
        join(SAVE_PATH, filename, "OUTPUT.png"),
      ];

      const p = Deno.run({ cmd: cmd });
      if ((await p.status()).code != 0) {
        throw Error("Bad image uploaded! Try another.");
      }
      p.close();
      jobQueue.add(filename);
      return new Response(`Success! Your ID is ${filename}`, {
        status: 201,
      });
    } catch (error) {
      console.error(error);
      return new Response(
        "Something went wrong! Please wait before trying again. " +
          error.message,
        { status: 500 },
      );
    }
  } else if (req.method == "DELETE") {
    try {
      const id = req.url.substring(req.url.lastIndexOf("/") + 1);
      Deno.remove(join(SAVE_PATH, id, "RESULT.png"));
      Deno.remove(join(SAVE_PATH, id), { recursive: true });
      return new Response("Content deleted.", { status: 200 });
    } catch (error) {
      console.error(error);
      return new Response("Could not delete content.", { status: 400 });
    }
  } else {
    return new Response(
      "Not a valid request.",
      { status: 400 },
    );
  }
}
