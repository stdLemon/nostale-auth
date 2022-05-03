#!/bin/node
import fs from "fs"
import {BlackboxEncoding} from "./blackbox_encoding.js"

const filename = process.argv[2]
const blackbox = fs.readFileSync(filename, {encoding: "utf8", flag: 'r'});
const fingerprint = BlackboxEncoding.decode(blackbox)
console.log(JSON.stringify(fingerprint, null, "\t"))

