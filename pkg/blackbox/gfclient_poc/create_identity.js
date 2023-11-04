#!/bin/node
import fs from "fs"
import {BlackboxEncoding} from "./blackbox_encoding.js"
import {Blackbox} from "./blackbox.js"
import {faker} from "@faker-js/faker"

const IDENTITY_FIELDS = ["v", "tz", "dnt", "product", "osType", "app", "vendor", "cookies", "mem", "con", "lang", "plugins", "gpu", "fonts", "audioC", "analyser", "width", "height", "depth", "video", "audio", "media", "permissions", "audioFP", "webglFP", "canvasFP", "osVersion", "userAgent"]

function main() {
    const filename = process.argv[2]
    const blackbox = fs.readFileSync(filename, {encoding: "utf8", flag: 'r'});

    const fingerprint = BlackboxEncoding.decode(blackbox)

    const identity = {
        timing: {min: 0, max: 0},
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
