import { normalize, basename, join, PQueue } from "./deps.ts";

const SAVE_PATH = Deno.env.get("SAVE_PATH") || "../receive/saves";
const QUEUE_PATH = Deno.env.get("QUEUE_PATH") || "../receive/queue";
const EXECUTABLE_PATH = Deno.env.get("EXECUTABLE_PATH") || "./thread-art-rust";

async function workOnID(id: string): Promise<void> {
  const inputPath = join(SAVE_PATH, id, "OUTPUT.png");
  const outputPath = join(SAVE_PATH, id, "RESULT.png");
  const cmd = [EXECUTABLE_PATH, inputPath, outputPath];
  const p = Deno.run({ cmd: cmd });
  if((await p.status()).code != 0) {
    const errorPath = join(SAVE_PATH, id, "ERROR")
    await Deno.create(errorPath);
    console.log("The algorithm crashed on", id + ".");
  } else {
    console.log("Done with", id + ".");
  }
  p.close();
}

const queue = new PQueue({ concurrency: 2 });

try {
  const queueDir = Deno.readDir(QUEUE_PATH);
  for await (const queueEntry of queueDir) {
    const queueID = queueEntry.name;
    queue.add(() => workOnID(queueID));
  }
} catch {
  // nothing is done
}

const watcher = Deno.watchFs(QUEUE_PATH);
for await (const event of watcher) {
  const filePath = normalize(event.paths[0]);
  const idFromFile = basename(filePath);
  try {
    if ("OK" == await Deno.readTextFile(filePath)) {
      await Deno.remove(filePath);
      queue.add(() => workOnID(idFromFile));
      console.log(idFromFile + " has been added to the queue.");
    }
  } catch (_error) {
    // nothing happens
  }
}
