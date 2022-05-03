import faker from '@faker-js/faker';
import axios from "axios"
import {Blackbox} from "./blackbox.js"

export class GfLauncher {
    constructor(gf_user_agent, cef_user_agent, installation_id) {
        this.cef_user_agent = cef_user_agent
        this.installation_id = installation_id
        this.bearer = null

        this.gf_headers = {
            "TNT-Installation-Id": installation_id,
            "Origin": "spark://www.gameforge.com",
            "User-Agent": gf_user_agent
        }
    }

    async auth(email, password, locale) {
        const url = "https://spark.gameforge.com/api/v1/auth/sessions"

        const data = {
            email,
            password,
            locale
        }
        const r = await axios.post(url, data, {headers: this.gf_headers})
        if (r.data.token !== undefined) {
            this.bearer = r.data.token
            return true
        }
        return false
    }

    async get_accounts() {
        const url = "https://spark.gameforge.com/api/v1/user/accounts"
        const headers = {
            ...this.gf_headers,
            "Authorization": `Bearer ${this.bearer}`,
        }

        const r = await axios.get(url, {headers})
        return r.data
    }

    async iovation(identity, accout_id) {
        const url = "https://spark.gameforge.com/api/v1/auth/iovation"
        const fingerprint = await Blackbox.create_fingerprint(identity)
        const blackbox = Blackbox.create_blackbox(fingerprint)

        const headers = {
            ...this.gf_headers,
            "Authorization": `Bearer ${this.bearer}`,
        }

        const data = {
            accoutId: accout_id,
            blackbox,
            type: "play_now"
        }

        const r = await axios.post(url, data, {headers})
        return r.data
    }

    async codes(identity, game_id, account_id) {
        const url = "https://spark.gameforge.com/api/v1/auth/thin/codes"
        const gs_id = this.#generate_gsid()
        const fingerprint = await Blackbox.create_fingerprint(identity)
        const encrypted_blackbox = Blackbox.create_encrypted_blackbox(fingerprint, gs_id, account_id, this.installation_id)

        const headers = {
            "tnt-installation-id": this.installation_id,
            "Authorization": `Bearer ${this.bearer}`,
            "User-Agent": this.cef_user_agent
        }

        const data = {
            blackbox: encrypted_blackbox,
            gameId: game_id,
            gsid: gs_id,
            platformGameAccountId: account_id
        }

        const r = await axios.post(url, data, {headers})
        return r.data
    }

    #generate_gsid() {
        const session = faker.datatype.uuid()
        const num = faker.datatype.number({min: 1, max: 9999}).toString()
        return `${session}-${num.padStart(4, "0")}`
    }
}
