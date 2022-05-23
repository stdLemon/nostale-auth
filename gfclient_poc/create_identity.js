#!/bin/node
import fs from "fs"
import {BlackboxEncoding} from "./blackbox_encoding.js"
import {Blackbox} from "./blackbox.js"
import {faker} from "@faker-js/faker"

const IDENTITY_FIELDS = ['v', 'tz', "dnt", "product", 'osType', 'app', "vendor", "cookies", 'mem', 'con', "lang", "plugins", 'gpu', "fonts", "audioC", "analyser", 'width', 'height', "depth", 'lStore', "sStore", "video", "audio", "media", "permissions", 'audioFP', 'webglFP', "canvasFP", "osVersion", 'userAgent']

function main() {
    const filename = process.argv[2]
    const blackbox = fs.readFileSync(filename, {encoding: "utf8", flag: 'r'});

    const fingerprint = BlackboxEncoding.decode(blackbox)

    const identity = {
        timing:
        {
            dP: {min: 0, max: 0}, // browser_info, platform_info, perms_media_audio timing
            dF: {min: 0, max: 0}, // fonts timing
            dW: {min: 0, max: 0}, // webgl timing
            dC: {min: 0, max: 0}, // canvas timing
            d: {min: 0, max: 0}, // sum of timings + few ms
        },
        fingerprint: {},
        installation_id: faker.datatype.uuid()
    }

    for (const field of IDENTITY_FIELDS) {
        identity.fingerprint[field] = fingerprint[field]
    }

    identity.fingerprint.vector = Blackbox.generate_vector()
    identity.fingerprint.uuid = Blackbox.generate_uuid()
    const indentity_json = JSON.stringify(identity, null, "\t")
    console.log(indentity_json)
}

main()
