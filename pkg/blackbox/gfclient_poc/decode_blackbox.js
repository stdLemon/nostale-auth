#!/bin/node
import fs from "fs";
import { BlackboxEncoding } from "./blackbox_encoding.js";

const filename = process.argv[2];
const blackbox = fs.readFileSync(filename, { encoding: "utf8", flag: "r" });
try {
    const fingerprint = BlackboxEncoding.decode(blackbox);
    console.log(JSON.stringify(fingerprint, null, "\t"));
} catch (e) {
    console.error("Failed to decode blackbox data:", e);
    process.exit(1);
}
