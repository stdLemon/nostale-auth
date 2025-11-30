#!/bin/node
import fs from "fs";
import { BlackboxEncoding } from "./blackbox_encoding.js";
import { BlackboxEncryption } from "./blackbox_encryption.js";

const filename = process.argv[2];
const code_request = JSON.parse(fs.readFileSync(filename, { encoding: "utf8", flag: "r" }));
const blackbox = BlackboxEncryption.decrypt(
    code_request.blackbox,
    code_request.gsid,
    code_request.platformGameAccountId,
);
const fingerprint = BlackboxEncoding.decode(blackbox);
console.log(JSON.stringify(fingerprint, null, "\t"));
