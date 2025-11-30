#!/bin/node
import fs from "fs";
import { BlackboxEncoding } from "./blackbox_encoding.js";

const filename = process.argv[2];
const fingerprint = JSON.parse(fs.readFileSync(filename, { encoding: "utf8", flag: "r" }));
const blackbox = BlackboxEncoding.encode(fingerprint);
process.stdout.write(blackbox);
