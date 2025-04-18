// This file contains svelte store implementations

import { writable } from "svelte/store";
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

export const jwt = writable<string>(local.accessToken || "");
jwt.subscribe(j => local.accessToken = j);
