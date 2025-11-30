#!/bin/node
import fs from "fs";
import { BlackboxEncryption } from "./blackbox_encryption.js";

const blackbox_filename = process.argv[2];
const gsid = process.argv[3];
const account_id = process.argv[4];

const blackbox = fs.readFileSync(blackbox_filename, { encoding: "utf8", flag: "r" });
const encrypted_blackbox = BlackboxEncryption.encrypt(blackbox, gsid, account_id);
process.stdout.write(encrypted_blackbox);
