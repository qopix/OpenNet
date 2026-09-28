import { connect } from "cloudflare:sockets";

export default {
  async fetch(request, env) {
    if (request.headers.get("Upgrade")?.toLowerCase() !== "websocket") {
      return new Response("OpenNet transport\n", { status: 426 });
    }

    const auth = request.headers.get("Authorization") || "";

    if (!env.OPENNET_TOKEN || auth !== `Bearer ${env.OPENNET_TOKEN}`) {
      return new Response("Unauthorized\n", { status: 401 });
    }

    const pair = new WebSocketPair();
    const client = pair[0];
    const server = pair[1];

    server.accept();

    let socket = null;
    let writer = null;
    let connected = false;

    const closeAll = () => {
      try {
        writer?.releaseLock();
      } catch {}

      try {
        socket?.close();
      } catch {}

      try {
        server.close();
      } catch {}
    };

    server.addEventListener("message", async (event) => {
      try {
        if (typeof event.data === "string") {
          if (!event.data.startsWith("CONNECT ")) {
            server.close(1002, "bad command");
            return;
          }

          const target = event.data.slice(8);
          const separator = target.lastIndexOf(":");

          if (separator <= 0) {
            server.close(1002, "bad target");
            return;
          }

          const hostname = target.slice(0, separator);
          const port = Number(target.slice(separator + 1));

          if (
            !hostname ||
            !Number.isInteger(port) ||
            port < 1 ||
            port > 65535
          ) {
            server.close(1002, "bad target");
            return;
          }

          socket = connect({
            hostname,
            port
          });

          writer = socket.writable.getWriter();
          connected = true;

          server.send("OK");

          const reader = socket.readable.getReader();

          try {
            while (true) {
              const { value, done } = await reader.read();

              if (done) break;

              if (value) {
                server.send(value);
              }
            }
          } finally {
            reader.releaseLock();
            closeAll();
          }

          return;
        }

        if (!connected || !writer) {
          server.close(1002, "not connected");
          return;
        }

        let data;

        if (event.data instanceof ArrayBuffer) {
          data = new Uint8Array(event.data);
        } else if (event.data instanceof Blob) {
          data = new Uint8Array(await event.data.arrayBuffer());
        } else {
          return;
        }

        await writer.write(data);
      } catch (error) {
        console.log(String(error));
        closeAll();
      }
    });

    server.addEventListener("close", closeAll);
    server.addEventListener("error", closeAll);

    return new Response(null, {
      status: 101,
      webSocket: client
    });
  }
};
