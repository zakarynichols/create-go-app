import { dirname } from "path";
import { fileURLToPath } from "url";
import express from "express";
import { join } from "path";
import { writeFileSync } from "fs";
import { indexHTML } from "./template.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = express();

const client = {
  port: process.env.NODE_CLIENT_PORT,
  host: process.env.NODE_CLIENT_HOST,
};

if (client.port === undefined) {
  throw new Error("missing port environment variable");
}

if (client.host === undefined) {
  throw new Error("missing host environment variable");
}

const data = new Uint8Array(Buffer.from(indexHTML));
writeFileSync("index.html", data);

app.use(express.static(join(__dirname)));

app.get("/", (_req, res) => {
  res.sendFile(join(__dirname, "../index.html"));
});

app.listen(client.port, () => {
  console.log(`Serving index.html on port: ${client.port}`);
});
