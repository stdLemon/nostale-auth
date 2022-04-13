#!/bin/node
const fs = require('fs');

const FIELD_LIST = ['v', 'tz', "dnt", "product", 'osType', 'app', "vendor", "cookies", 'mem', 'con', "lang", "plugins", 'gpu', "fonts", "audioC", "analyser", 'width', 'height', "depth", 'lStore', "sStore", "video", "audio", "media", "permissions", 'audioFP', 'webglFP', "canvasFP", 'dP', 'dF', 'dW', 'dC', "creation", "uuid", 'd', "osVersion", "vector", 'userAgent', "serverTimeInMS", "request"]

function decode_blackbox(blackbox) {
    const gf_encoded = atob(blackbox.replaceAll("_", "/").replaceAll("-", "+"))
    let uri_encoded = gf_encoded[0]
    for (let i = 1; i < gf_encoded.length; ++i) {
        const b = gf_encoded.charCodeAt(i - 1)
        let a = gf_encoded.charCodeAt(i)

        if (a < b) {
            a += 0x100
        }

        const c = String.fromCharCode(a - b)
        uri_encoded += c
    }
    const fingerprint = decodeURIComponent(uri_encoded)
    return fingerprint
}

const filename = process.argv[2]
const blackbox = fs.readFileSync(filename, {encoding: "utf8", flag: 'r'});


const fingerprint = JSON.parse(decode_blackbox(blackbox))

fingerprint_obj = {}
for (i in FIELD_LIST) {
    fingerprint_obj[FIELD_LIST[i]] = fingerprint[i]
}
console.log(fingerprint_obj)

