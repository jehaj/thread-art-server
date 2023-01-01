import { join } from "https://deno.land/std@0.170.0/path/posix.ts";
import { serve } from "https://deno.land/std@0.170.0/http/server.ts";

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

const QUEUE_PATH = "./queue";
const QUEUE_LIMIT = 10;

serve(handler, { port: 8001 });

async function handler(req: Request): Promise<Response> {
  const d = Deno.readDir("./queue");
  let sum = 0;
  for await (const _ of d) {
    sum += 1;
  }
  if (sum >= QUEUE_LIMIT) {
    return new Response("Try again later!", { status: 500 });
  }
  try {
    const f = await req.formData();
    const imageEntry = f.get("image");
    if (imageEntry == null || typeof imageEntry != "object") {
      throw new Error("Image not contained in form data.");
    }
    const image: Uint8Array = new Uint8Array(await imageEntry.arrayBuffer());
    const filename = RandomID();
    await Deno.mkdir(`./saves/${filename}`);
    await Deno.writeFile(`./saves/${filename}/INPUT.png`, image, {
      createNew: true,
    });
    const cmd = [
      "convert",
      `./saves/${filename}/INPUT.png`,
      "-resize",
      "400x400^",
      "-gravity",
      "center",
      "-extent",
      "400x400",
      "-colorspace",
      "Gray",
      `./saves/${filename}/OUTPUT.png`,
    ];
    
    const p = Deno.run({ cmd: cmd });
    await p.status();
    await Deno.writeTextFile(join(QUEUE_PATH, filename), "");
    p.close();
    
    return new Response(`Success! Your ID is ${filename}`, { status: 200 });
    
  } catch (error) {
    console.log(error);
    return new Response(
      "Something went wrong! Please wait before trying again. " + error.message,
      { status: 500 },
    );
  }
}
