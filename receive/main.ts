const server = Deno.listen({ port: 8001 });
for await (const conn of server) {
  serveHttp(conn);
}

async function serveHttp(conn: Deno.Conn) {
  const httpConn = Deno.serveHttp(conn);
  for await (const requestEvent of httpConn) {
    const f = await requestEvent.request.formData();
    const image = await f.get("image").arrayBuffer();
    console.log(image);
    await Deno.writeFile('./saves/image.png', image);
    requestEvent.respondWith(
      new Response("Success!", { status: 200 })
    );
  }
}
