import { Job, join, Worker } from "./deps.ts";

const SAVE_PATH = Deno.env.get("SAVE_PATH") || "../receive/saves";
const EXECUTABLE_PATH = Deno.env.get("EXECUTABLE_PATH") || "./thread-art-rust";
const CONCURRENCY = parseInt(Deno.env.get("CONCURRENCY") || "2");

const worker = new Worker("job", workOnID, {
  concurrency: CONCURRENCY,
  connection: {
    host: "localhost",
  },
});

console.log("Act has started.");

async function workOnID(job: Job): Promise<void> {
  const id = job.name;
  const inputPath = join(SAVE_PATH, id, "OUTPUT.png");
  const outputPath = join(SAVE_PATH, id, "RESULT.png");
  const cmd = [EXECUTABLE_PATH, inputPath, outputPath];
  const p = Deno.run({ cmd: cmd });
  if ((await p.status()).code != 0) {
    const errorPath = join(SAVE_PATH, id, "ERROR");
    await Deno.create(errorPath);
    console.log("The algorithm crashed on", id + ".");
  } else {
    console.log("BullMQ: Done with", id + ".");
  }
  p.close();
}
