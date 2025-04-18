// This file contains helper functions for user auth

import { get as value , type Writable } from 'svelte/store'
import page from 'page'
import request from "./request";
import { jwt, user } from "./stores";
import toast from "./toast";
import type { JwtUserData } from "./types";

// Permission levels for a user
export enum Permission {
    None = 0,
    Read = 4,
    Write = 2,
    Admin = 1
}

/**
 * Check if there is an active user
 */
export const isLoggedIn = (user: JwtUserData | null): boolean => {
    return user != null;
};

/**
 * Check if the user in question has access to a resource
 */
export const hasAccess = (user: JwtUserData|null, level: Permission = Permission.None): boolean => {
    if (!user) {
        return false
    }

    return (user.level & level) === level
}

/**
 * Guard route access to guest only
 */
export const guardGuest = (_, next: () => void) => {
    if (isLoggedIn(value(user))) {
        page("/home")
        return
    }

    next()
}

/**
 * Guard route access to logged in users only
 */
export const guardUser = (_, next: () => void) => {
    if (!isLoggedIn(value(user))) {
        page("/account/login")
        return
    }

    next()
}

/**
 * Guard route access to users with read access only
 */
export const guardReadUser = (_, next: () => void) => {
    if (!isLoggedIn(value(user))) {
        page("/account/login")
        return
    }

    if (!hasAccess(value(user), Permission.Read)) {
        toast.error("Access Denied")
        page("/home")
        return
    }

    next()
}

/**
 * Guard route access to users with write access only
 */
export const guardWriteUser = (_, next: () => void) => {
    if (!isLoggedIn(value(user))) {
        page("/account/login")
        return
    }

    if (!hasAccess(value(user), Permission.Write)) {
        toast.error("Access Denied")
        page("/home")
        return
    }

    next()
}

/**
 * Guard route access to users with admin access only
 */
export const guardAdminUser = (_, next: () => void) => {
    if (!isLoggedIn(value(user))) {
        page("/account/login")
        return
    }

    if (!hasAccess(value(user), Permission.Admin)) {
        toast.error("Access Denied")
        page("/home")
        return
    }

    next()
}

/**
 * Attempt to verify the sesison
 */
export const verifySession = async (): Promise<void> => {
    try {
        await request.get("/api/auth/verify");
        if (window.location.pathname == "/") {
            console.log("redirecting")
            page.redirect("/home")
        }
    } catch (e) {
        page("/account/login")
    }
};

/**
 * Attempt to login with the sso provider
 */
export const login = async (user: string, passwd: string): Promise<void> => {
    try {
        const resp =  await request.post("/api/auth/login", { user, passwd })
        console.log({ resp })

        toast.alert("Logged in");
        page("/home");
    } catch (e) {
        console.log(e)
        toast.error("Login failed");
    }
};

/**
 * Attempt to logout with the sso provider
 */
export const logout = async (): Promise<void> => {
    try {
        jwt.set("");
        user.set(null);

        toast.alert("Logged out");
        page("/account/login");
    } catch (e) {
        console.log(e)
        toast.error("Logout failed");
    }
};
