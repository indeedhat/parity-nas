// This file contains helper functions for managing JWT's

import type { JwtToken } from "./types";

/**
 * Parse a jwt token from string into a usable object
 */
export const parse = (token: string): JwtToken | null => {
    const body64 = token.split(".")[1];
    const normal64 = body64.replace(/-/g, '+').replace(/_/g, '/');
    const json = decodeURIComponent(
        window.atob(normal64)
            .split('')
            .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
            .join('')
    );

    return JSON.parse(json);
};

export default {
    parse
};
