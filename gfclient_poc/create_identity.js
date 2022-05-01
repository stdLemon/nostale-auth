#!/bin/node
import fs from "fs"
import GfclientEncoding from "./gfclient_encoding.js"

const IDENTITY_FIELDS = ['v', 'tz', "dnt", "product", 'osType', 'app', "vendor", "cookies", 'mem', 'con', "lang", "plugins", 'gpu', "fonts", "audioC", "analyser", 'width', 'height', "depth", 'lStore', "sStore", "video", "audio", "media", "permissions", 'audioFP', 'webglFP', "canvasFP", "uuid", "osVersion", "vector", 'userAgent', "request"]

function main() {
    const filename = process.argv[2]
    const blackbox = fs.readFileSync(filename, {encoding: "utf8", flag: 'r'});

    const fingerprint = GfclientEncoding.decode_blackbox(blackbox)

    const identity = {
        timing:
        {
            dP: {min: 0, max: 0}, // browser_info, platform_info, perms_media_audio timing
            dF: {min: 0, max: 0}, // fonts timing
            dW: {min: 0, max: 0}, // webgl timing
            dC: {min: 0, max: 0}, // canvas timing
            d: {min: 0, max: 0}, // sum of timings + few ms
        },
        fingerprint: {}
    }

    for (let field of IDENTITY_FIELDS) {
        identity.fingerprint[field] = fingerprint[field]
    }

    identity.fingerprint.vector = Buffer.from(fingerprint.vector, "base64").toString("latin1")
    const indentity_json = JSON.stringify(identity, null, "\t")
    console.log(indentity_json)
}

main()
