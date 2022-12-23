function RandomID() {
  const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const clen = letters.length;
  const r = crypto.getRandomValues(new Uint8Array(12));
  const b = [];
  
  for (let i = 0; i < 12; i++) {
    const c = parseInt(r[i]);
    b.push(letters[c%clen]);
  }
  
  return b.join('');
}

const server = Deno.listen({ port: 8001 });
for await (const conn of server) {
  serveHttp(conn);
}

async function serveHttp(conn: Deno.Conn) {
  const httpConn = Deno.serveHttp(conn);
  for await (const requestEvent of httpConn) {
    try {
      const f = await requestEvent.request.formData();
      const image = await f.get("image").arrayBuffer();
      const filename = RandomID();
      await Deno.writeFile(`./saves/${filename}.png`, image);
      requestEvent.respondWith(
        new Response(`Success! Your ID is ${filename}`, { status: 200 })
      );
    } catch (error) {
      requestEvent.respondWith(
        new Response("Something went wrong! Please wait before trying again.", { status: 500 })
      );
    }
  }
}
