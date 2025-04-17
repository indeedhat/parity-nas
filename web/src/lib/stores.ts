// This file contains svelte store implementations

import { writable,  get as value  } from "svelte/store";
import type { JwtUserData } from "./types";

const local = "undefined" != typeof localStorage
    ? localStorage
    : { userData: null, accessToken: "" };

export const user = writable<JwtUserData | null>(
    local.userData
        ? JSON.parse(local.userData)
        : null
);
user.subscribe(usr => local.userData = JSON.stringify(usr));

console.log(local, local.accessToken || "")
export const jwt = writable<string>(local.accessToken || "");
console.log(value(jwt))
jwt.subscribe(j => local.accessToken = j);
