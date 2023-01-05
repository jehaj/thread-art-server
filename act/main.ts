import { normalize, basename } from "https://deno.land/std@0.170.0/path/mod.ts";
import { join } from "https://deno.land/std@0.170.0/path/mod.ts";
import PQueue from "https://deno.land/x/p_queue@1.0.1/mod.ts";

const SAVE_PATH = Deno.env.get("SAVE_PATH") || "../receive/saves";
const QUEUE_PATH = Deno.env.get("QUEUE_PATH") || "../receive/queue";

async function workOnID(id: string): Promise<void> {
  const inputPath = join(SAVE_PATH, id, "OUTPUT.png");
  const outputPath = join(SAVE_PATH, id, "RESULT.png");
  const cmd = ["./thread-art-rust", inputPath, outputPath];
  const p = Deno.run({ cmd: cmd });
  await p.status();
  console.log("Done with", id + ".");
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
