import faker from "@faker-js/faker"
import axios from "axios"
import {BlackboxEncoding} from "./blackbox_encoding.js"
import {BlackboxEncryption} from "./blackbox_encryption.js"

const SERVER_FILE_GAME1_FILE = "https://gameforge.com/tra/game1.js"
const VECTOR_LENGTH = 100

function random_ascii_character() {
    return String.fromCharCode(0x20 + Math.random() * (0x7E - 0x20) | 0x0)
}

function generate_vector() {
    const vector_content = Array.from(Array(VECTOR_LENGTH), random_ascii_character).join('')
    const time = new Date().getTime()
    return `${vector_content} ${time}`
}

function generate_uuid(arg = 3) {
    return new Array(arg).fill(0x0).map(() => Math.random().toString(0x24).substr(0x2, 0x9)).reduce((a, b) => a + b, '');
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

function create_blackbox(fingerprint, request = null) {
    fingerprint.request = request
    return BlackboxEncoding.encode(fingerprint)
}

function create_encrypted_blackbox(fingerprint, gs_id, account_id, installation) {
    const delim_index = gs_id.lastIndexOf("-")
    const session = gs_id.substring(0, delim_index)

    const request = {
        features: [faker.datatype.number({min: 1, max: 0xFFFFFFFE - 1})],
        installation,
        session
    }

    const blackbox = create_blackbox(fingerprint, request)
    return BlackboxEncryption.encrypt(blackbox, gs_id, account_id)
}

export const Blackbox = {
    create_fingerprint,
    create_blackbox,
    create_encrypted_blackbox,
    generate_uuid,
    generate_vector,
    update_vector
}
