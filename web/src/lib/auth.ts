// This file contains helper functions for user auth

import page from 'page'
import request from "./request";
import { jwt, user } from "./stores";
import toast from "./toast";
import type { JwtUserData } from "./types";

/**
 * Check if there is an active user
 */
export const isLoggedIn = (user: JwtUserData | null): boolean => {
    return user != null;
};

/**
 * Check if there is an active user
 */
export const redirectGuests = (user: JwtUserData | null): void => {
    if (isLoggedIn(user)) {
        return
    }

    page("/account/login")
};

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
