#!/bin/node
import faker from "@faker-js/faker"
import fs from "fs"
import axios from "axios"
import GfclientEncoding from "./gfclient_encoding.js"

const SERVER_FILE_GAME1_FILE = "https://gameforge.com/tra/game1.js"

function random_ascii_character() {
    return String.fromCharCode(0x20 + Math.random() * (0x7E - 0x20) | 0x0)
}

async function get_server_date() {
    const res = await axios.get(SERVER_FILE_GAME1_FILE)
    return new Date(res.headers["date"]).toISOString()
}

function update_vector(vector) {
    const delim_index = vector.lastIndexOf(" ")
    let vec_content = vector.substring(0, delim_index)
    let vec_time = vector.substring(delim_index + 1)
    const current_time = new Date().getTime()

    vec_content = vec_content.split('')
    vec_time = parseInt(vec_time)

    if (current_time > vec_time + 1000) {
        vec_content.shift()
        vec_content.push(random_ascii_character())

        const new_vec = `${vec_content.join("")} ${current_time}`
        return new_vec
    }

    return vector
}

async function create_fingerprint(identity) {
    const fingerprint = {
        ...identity.fingerprint,
        dP: faker.datatype.number(identity.timing.dP), // timing browser_info, platform_info, perms_media_audio 
        dF: faker.datatype.number(identity.timing.dF), // intalled fonts fignerprint timing
        dW: faker.datatype.number(identity.timing.dW), // webgl fingerprint timing
        dC: faker.datatype.number(identity.timing.dC), // canvas fignerprint timing
        creation: new Date().toISOString(),
        serverTimeInMS: await get_server_date(),
    }
    fingerprint.d = fingerprint.dP + fingerprint.dF + fingerprint.dW + fingerprint.dC + faker.datatype.number(identity.timing.d)
    fingerprint.vector = Buffer.from(fingerprint.vector).toString("base64")
    return fingerprint
}

async function main() {
    const filename = process.argv[2]
    const identity = JSON.parse(fs.readFileSync(filename, {encoding: "utf8", flag: 'r'}))

    identity.fingerprint.vector = update_vector(identity.fingerprint.vector)
    console.log("identity", identity)
    console.log()
    fs.writeFileSync(filename, JSON.stringify(identity, null, "\t"))
    const fingerprint = await create_fingerprint(identity)
    console.log("fingerprint", fingerprint)
    console.log()
    const blackbox = GfclientEncoding.encode_blackbox(fingerprint)
    console.log("blackbox", blackbox)
}

main()
