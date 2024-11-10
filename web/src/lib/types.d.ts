// This file contains all the custom types for the ui

// Dictionary with string values
export interface StringDict {
    [key: string]: string;
}

export interface Dict {
    [key: string]: any
}

export interface ToastNotification {
    id: ?number;
    message: string;
    type: ?string;
    dismissible: ?boolean;
    timeout: ?number,
}

export interface JwtToken {
    nme: string;
    uid: string;
    iss: string;
    sub: number;
    iat: number;
    exp: number;
    vbf: number;
    cct: number;
    lwy: number;
}

export interface JwtUserData {
    id: string
    name: string
}
