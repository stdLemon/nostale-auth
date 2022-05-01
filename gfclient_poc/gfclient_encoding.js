const BLACKBOX_FIELDS = ['v', 'tz', "dnt", "product", 'osType', 'app', "vendor", "cookies", 'mem', 'con', "lang", "plugins", 'gpu', "fonts", "audioC", "analyser", 'width', 'height', "depth", 'lStore', "sStore", "video", "audio", "media", "permissions", 'audioFP', 'webglFP', "canvasFP", 'dP', 'dF', 'dW', 'dC', "creation", "uuid", 'd', "osVersion", "vector", 'userAgent', "serverTimeInMS", "request"]

function encode_blackbox(fingerprint) {
    const fingerprint_array = []
    for (const field of BLACKBOX_FIELDS) {
        fingerprint_array.push(fingerprint[field])
    }

    const uri_encoded = encodeURIComponent(JSON.stringify(fingerprint_array))
    let gf_encoded = uri_encoded[0]

    for (let i = 1; i < uri_encoded.length; ++i) {
        const a = gf_encoded.charCodeAt(i - 1)
        const b = uri_encoded.charCodeAt(i)
        const c = String.fromCharCode((a + b) % 0x100);

        gf_encoded += c
    }

    return "tra:" + Buffer.from(gf_encoded).toString("base64")
        .replaceAll("/", "_")
        .replaceAll("+", "-")
        .replaceAll("=", "")
}

function decode_blackbox(blackbox, as_object = true) {

    blackbox = blackbox
        .replaceAll("tra:", "")
        .replaceAll("_", "/")
        .replaceAll("-", "+")
    const gf_encoded = Buffer.from(blackbox, "base64")
    let uri_encoded = String.fromCharCode(gf_encoded[0])

    for (let i = 1; i < gf_encoded.length; ++i) {
        const b = gf_encoded[i - 1]
        let a = gf_encoded[i]

        if (a < b) {
            a += 0x100
        }

        const c = String.fromCharCode(a - b)
        uri_encoded += c
    }

    const fingerprint_str = decodeURIComponent(uri_encoded.toString("latin1"))
    if (!as_object) {
        return fingerprint_str
    }

    const fingerprint_array = JSON.parse(fingerprint_str)
    const fingerprint = {}

    for (let i in BLACKBOX_FIELDS) {
        fingerprint[BLACKBOX_FIELDS[i]] = fingerprint_array[i]
    }

    return fingerprint
}

export default {
    encode_blackbox,
    decode_blackbox
}
