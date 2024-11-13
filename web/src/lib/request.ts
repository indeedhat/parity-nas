// This file contains a simple wrapper around the fetch api for more readable http requests

import { get as value } from 'svelte/store';
import { jwt, user } from '$lib/stores'
import { parse } from '$lib/jwt';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { WebRoot } from '$lib/env';
import type { Dict, StringDict } from './types';


/**
 * Perform an Http GET request
 */
export const get = async (url: string): Promise<unknown> => {
    return noBodyRequest("GET", url);
};

/**
 * Perform an Http DELETE request
 */
export const Delete = async (url: string): Promise<unknown> => {
    return noBodyRequest("DELETE", url);
};

/**
 * Perform an Http PUT request
 */
export const put = async (url: string, data: FormData | Dict, optHeaders: StringDict | null = null): Promise<unknown> => {
    return bodyRequest("PUT", url, data, optHeaders);
};

/**
 * Perform an Http POST request
 */
export const post = async (url: string, data: FormData | Dict, optHeaders: StringDict | null = null): Promise<unknown> => {
    return bodyRequest("POST", url, data, optHeaders);
};

const noBodyRequest = async (method: string, url: string): Promise<unknown> => {
    let opts = undefined;
    if (value(jwt)) {
        opts = <RequestInit>{
            method,
            credentials: 'include',
            headers: new Headers({
                'Authorization': `Bearer ${value(jwt)}`
            })
        };
    }

    const resp = await fetch(WebRoot + url, opts);

    handleAuth(resp);

    if (resp.headers.get('Content-Type') === 'application/json') {
        return await resp.json();
    }

    return await resp.text()
}

const bodyRequest = async (
    method: string,
    url: string,
    data: FormData | Dict,
    optHeaders: StringDict | null = null

): Promise<unknown> => {
    const headers = new Headers();

    if (optHeaders) {
        for (const k in optHeaders) {
            headers.set(k, optHeaders[k])
        }
    }
    if (value(jwt)) {
        headers.set("Authorization", `Bearer ${value(jwt)}`);
    }

    const resp = await fetch(WebRoot + url, <RequestInit>{
        method,
        body: object2form(data),
        credentials: 'include',
        headers: value(jwt) ? new Headers({
            'Authorization': `Bearer ${value(jwt)}`
        }) : undefined
    })

    handleAuth(resp);

    if (resp.headers.get('Content-Type') === 'application/json') {
        return await resp.json();
    }

    return await resp.text()
}

const object2form = (obj: Dict | FormData) => {
    if (obj instanceof FormData) {
        return obj;
    }

    return JSON.stringify(obj)
};

const handleAuth = (resp: Response) => {
    if (resp.status == 401) {
        jwt.set("")
        user.set(null);

        if (browser) {
            goto("/");
        }

        throw new NotAuthorized();
    }

    let token = resp.headers.get("Auth_token");
    resp.headers.forEach(function() {
        console.log(arguments)
    })
    console.log(token)
    if (token !== null && token.startsWith("jwt.")) {
        token = token.substring(4);
        console.log(token)
        jwt.set(token);

        const data = parse(token);
        if (data) {
            user.set({ name: data.nme, id: data.uid });
        }
    }
};

export class NotAuthorized {}


export default {
    get,
    post,
    put,
    Delete
};
