// This file provides access to environment variables
// The environment variables are provided by the server to the svelte app

import { browser } from "$app/environment";

let WebRoot = "http://localhost:8080";

if (browser && window.proc?.env?.loaded) {
    WebRoot = window.proc.env.root;
}

export default {
    WebRoot,
};

export {
    WebRoot,
};
