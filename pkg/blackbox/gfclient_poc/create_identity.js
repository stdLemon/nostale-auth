#!/bin/node
import fs from "fs";
import { randomUUID } from "node:crypto";
import { BlackboxEncoding } from "./blackbox_encoding.js";
import { Blackbox } from "./blackbox.js";

function main() {
    const filename = process.argv[2];
    const blackbox = fs.readFileSync(filename, { encoding: "utf8", flag: "r" });

    const identity = {
        timing: { min: 0, max: 0 },
        fingerprint: BlackboxEncoding.decode(blackbox),
        installation_id: randomUUID(),
    };

    identity.fingerprint.vecSignatureBase64 = Blackbox.generate_vector();
    identity.fingerprint.clientId = Blackbox.generate_client_id();
    const indentity_json = JSON.stringify(identity, null, "\t");
    console.log(indentity_json);
}

main();
