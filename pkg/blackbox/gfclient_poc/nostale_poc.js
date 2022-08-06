#!/bin/node
import fs from "fs"
import {GfLauncher} from "./gf_launcher.js"
import {Blackbox} from "./blackbox.js"

async function main() {
    const identity_path = "identity.json"
    const account_path = "account.json"

    const cef_user_agent = "Chrome/C2.2.23.1813 (49c0acbee1)"

    const identity = JSON.parse(fs.readFileSync(identity_path))
    const account_data = JSON.parse(fs.readFileSync(account_path))

    const gf_launcher = new GfLauncher(identity.fingerprint.userAgent, cef_user_agent, identity.installation_id)
    const auth_ok = await gf_launcher.auth(account_data.email, account_data.password, account_data.locale)
    if (!auth_ok) {
        console.error("Authorization failed")
        return
    }

    const accounts = Object.values(await gf_launcher.get_accounts())
    const account = accounts.find(acc => acc.displayName === account_data.name)
    if (account === undefined) {
        console.error("account with name", account_data.name, "was not found")
    }

    identity.fingerprint.vector = Blackbox.update_vector(identity.fingerprint.vector)
    const iovation_response = await gf_launcher.iovation(identity, account.id)
    console.log(iovation_response)

    identity.fingerprint.vector = Blackbox.update_vector(identity.fingerprint.vector)
    const codes_response = await gf_launcher.codes(identity, account.gameId, account.id)
    console.log(codes_response)

    fs.writeFileSync(identity_path, JSON.stringify(identity, null, "\t"))
}
main()
