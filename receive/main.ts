import { join } from "https://deno.land/std@0.170.0/path/mod.ts";
import { serve } from "https://deno.land/std@0.170.0/http/server.ts";
import "https://deno.land/std@0.170.0/dotenv/load.ts";

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
const QUEUE_PATH = Deno.env.get("QUEUE_PATH") || "./queue";
const PORT = parseInt(Deno.env.get("PORT") || "8001")

console.log("Running with settings:");
console.log("QUEUE_LIMIT =", QUEUE_LIMIT);
console.log("SAVE_PATH =", SAVE_PATH);
console.log("QUEUE_PATH =", QUEUE_PATH);
console.log("PORT =", PORT);
console.log("Working at:", Deno.cwd());

await Deno.mkdir(SAVE_PATH, { recursive: true });
await Deno.mkdir(QUEUE_PATH, { recursive: true });

serve(handler, { port: PORT });

async function handler(req: Request): Promise<Response> {
  const filename = RandomID();
  await Deno.writeTextFile(join(QUEUE_PATH, filename), "",);

  const d = Deno.readDir(QUEUE_PATH);
  let sum = 0;
  for await (const _ of d) {
    sum += 1;
  }
  if (sum > QUEUE_LIMIT) {
    await Deno.remove(join(QUEUE_PATH, filename));
    return new Response("Try again later!", { status: 500 });
  }
  try {
    const f = await req.formData();
    const imageEntry = f.get("image");
    if (imageEntry == null || typeof imageEntry != "object") {
      throw new Error("Image not contained in form data.");
    }
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
    await p.status();
    p.close();

    await Deno.writeTextFile(join(QUEUE_PATH, filename), "OK",);
    return new Response(`Success! Your ID is ${filename}`, { status: 200 });
  } catch (error) {
    console.error(error);
    try {
      await Deno.remove(join(QUEUE_PATH, filename));
    } catch (_error) {
      // nothing is done
    }

    return new Response(
      "Something went wrong! Please wait before trying again. " + error.message,
      { status: 500 },
    );
  }
}
