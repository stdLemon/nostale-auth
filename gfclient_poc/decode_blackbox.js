#!/bin/node
import fs from "fs"
import GfclientEncoding from "./gfclient_encoding.js"

const filename = process.argv[2]
const blackbox = fs.readFileSync(filename, {encoding: "utf8", flag: 'r'});


const fingerprint = GfclientEncoding.decode_blackbox(blackbox)
console.log(JSON.stringify(fingerprint, null, "\t"))

