#!/bin/node
const {faker} = require('@faker-js/faker');
const VERSION = 7
const VECTOR_LENGTH = 100

function encode_blackbox(fingerprint) {
    const uri_encoded = encodeURIComponent(fingerprint)
    let gf_encoded = uri_encoded[0]

    for (let i = 1; i < uri_encoded.length; ++i) {
        const a = gf_encoded.charCodeAt(i - 1)
        const b = uri_encoded.charCodeAt(i)
        const c = String.fromCharCode((a + b) % 0x100);

        gf_encoded += c
    }
    const blackbox = btoa(gf_encoded).replaceAll("/", "_").replaceAll("+", "-").replaceAll("=", "")
    return blackbox

}

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


function random_ascii_character() {
    return String.fromCharCode(0x20 + Math.random() * (0x7E - 0x20) | 0x0)
}

function generate_vector() {
    const vector = Array.from(Array(VECTOR_LENGTH), random_ascii_character).join('')
    const time = new Date().getTime()
    return btoa(`${vector} ${time}`)
}

function generate_uuid(arg = 3) {
    return new Array(arg).fill(0x0).map(() => Math.random().toString(0x24).substr(0x2, 0x9)).reduce((a, b) => a + b, '');
}

const BROWSER_ENGINE_MAPPING = {
    'Chrome': "Blink",
    'Opera': "Blink",
    'Mozilla': "Gecko",
    'Edge': "Blink",
    'Safari': 'WebKit',
    'Firefox': "Gecko",
    'Internet Explorer': "Trident"
}


const fingerprint = {
    v: VERSION,
    tz: faker.address.timeZone(), // Intl["DateTimeFormat"]()["resolvedOptions"]()["timeZone"]
    dnt: false, // navigator["doNotTrack"] || false
    product: 'Blink', // browser engine, user agent 
    osType: 'Windows', // os name, user agent
    app: 'Chrome', // browser name, user agent
    vendor: 'Google Inc.', // navigator["vendor"]
    cookies: true, // navigator["cookieEnabled"]
    mem: faker.datatype.number({min: 1, max: 64}), // navigator["deviceMemory"]
    con: faker.datatype.number({min: 2, max: 16}), // navigator["hardwareConcurrency"]
    lang: 'en-US,en', // navigator["languages"].join(',')
    plugins: '4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945',
    gpu: 'Google Inc.,Google SwiftShader',
    fonts: 'edf79a33390a558a40ab0992b1e276770bf29d77842c90a6bbf590e61fb326ef',
    audioC: 'd9af7aa1d00f202e8291fe49b9344f69746635eea53e7eace68c10f302cc933a',
    analyser: 'fb2d753714f0b65738b5a3324c616654674fdce7c0493523e88bff731603a648',
    width: 1920, // window["screen"]["availWidth"]
    height: 1040, // window.screen.availHeight
    depth: 24, // window["screen"]["colorDepth"]
    lStore: true, // Boolean(localStorage)
    sStore: true, // Boolean(sessionStorage)
    video: 'ea2c39c5eca488bd7ee0a1d7ce6b5600da5f36a7c8aa89bdeb078690fe8950e6',
    audio: '456687e4e0029125c0b73edf391fa02e5b8906ef5bc3ac2b34ebe93ba04c130a',
    media: 'c9b87e00ae04b169e849cbc814b70e626142fa9f37f792689957edaf8cead76a',
    permissions: '27f243aa0ac84a576f3009806c3b13614a4efb01a9a42420041da0b72ec4c9b6',
    audioFP: 124.0434474653739,
    webglFP: '301740ff533ffb146abf0cce712f56cb8d06ba6bb399cd571c092b5519289cb8',
    canvasFP: 1677023211,
    dP: faker.datatype.number({min: 40, max: 60}), // timing browser_info, platform_info, perms_media_audio 
    dF: faker.datatype.number({min: 155, max: 165}), // intalled fonts fignerprint timing
    dW: faker.datatype.number({min: 370, max: 420}), // webgl fingerprint timing
    dC: faker.datatype.number({min: 2, max: 10}), // canvas fignerprint timing
    creation: new Date().toISOString(),
    uuid: generate_uuid(),
    d: undefined, // sum of timings + ~15ms
    osVersion: '10', // optional
    vector: generate_vector(),
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36',
    serverTimeInMS: undefined, // same as creation time but ms zeroed out
    request: null
}
fingerprint.d = fingerprint.dP + fingerprint.dF + fingerprint.dW + fingerprint.dC + faker.datatype.number({min: 10, max: 15})
const server_time = new Date(fingerprint.creation)
server_time.setMilliseconds(0)
fingerprint.serverTimeInMS = server_time.toISOString()
const fingerprint_array = Object.values(fingerprint)
const fingerprint_json = JSON.stringify(fingerprint_array)
const blackbox = encode_blackbox(fingerprint_json)

const decoded_blackbox = decode_blackbox(blackbox)
console.log("fingerprint", fingerprint)
console.log("blackbox", blackbox)
console.log(fingerprint_json === decoded_blackbox)
