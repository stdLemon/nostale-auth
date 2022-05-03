import crypto from "crypto"

function xor(data, key) {
    const key_length = key.length
    const result = []

    for (let i = 0; i < data.length; ++i) {
        const wrapping_i = i % key_length

        const c = data[i] ^ key[wrapping_i] ^ key[key_length - wrapping_i - 1]
        result.push(c)
    }

    return Buffer.from(result)
}

function create_key(gs_id, account_id) {
    const hash = crypto.createHash("sha512")
    return Buffer.from(hash.update(`${gs_id}-${account_id}`, "utf8").digest("hex"))
}

function decrypt(encrypted_blackbox, gs_id, account_id) {
    return xor(
        Buffer.from(encrypted_blackbox, "base64"),
        create_key(gs_id, account_id)
    ).toString()
}

function encrypt(blackbox, gs_id, account_id) {
    return xor(
        Buffer.from(blackbox),
        create_key(gs_id, account_id)
    ).toString("base64")
}

export const BlackboxEncryption = {
    encrypt,
    decrypt
}
