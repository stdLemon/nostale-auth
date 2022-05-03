#!/bin/node
import fs from "fs"
import {GfLauncher} from "./gf_launcher.js"
import {Blackbox} from "./blackbox.js"

async function main() {
    const identity_path = "identity.json"
    const account_path = "account.json"
    const cef_user_agent = "Chrome/C2.2.23.1813 (49c0acbee1)"
    const account_index = 1

    const identity = JSON.parse(fs.readFileSync(identity_path))
    const account = JSON.parse(fs.readFileSync(account_path))

    const gf_launcher = new GfLauncher(identity.fingerprint.userAgent, cef_user_agent, account.installation_id)
    const auth_ok = await gf_launcher.auth(account.email, account.password, account.locale)

    if (!auth_ok) {
        console.log("Authorization failed")
        return
    }

    const accounts = await gf_launcher.get_accounts()
    const account_id = Object.keys(accounts)[account_index]
    const game_id = accounts[account_id].gameId

    identity.fingerprint.vector = Blackbox.update_vector(identity.fingerprint.vector)
    const iovation_response = await gf_launcher.iovation(identity, account_id)
    console.log(iovation_response)

    identity.fingerprint.vector = Blackbox.update_vector(identity.fingerprint.vector)
    const codes_response = await gf_launcher.codes(identity, game_id, account_id)
    console.log(codes_response)

    fs.writeFileSync(identity_path, JSON.stringify(identity, null, "\t"))
}
main()
