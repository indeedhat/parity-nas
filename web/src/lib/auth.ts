// This file contains helper functions for user auth

import { browser } from "$app/environment";
import { goto } from "$app/navigation";
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
 * Attempt to verify the sesison
 */
export const verifySession = async (): Promise<void> => {
    try {
        await request.get("/api/auth/verify");
    } catch (e) {
        if (browser) {
            goto("/")
        }
    }
};

/**
 * Attempt to login with the sso provider
 */
export const login = async (user: string, passwd: string): Promise<void> => {
    try {
        await request.post("/api/auth/login", { user, passwd })

        toast.alert("Logged in");
        if (browser) {
            goto("/home");
        }
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
        if (browser) {
            goto("/");
        }
    } catch (e) {
        console.log(e)
        toast.error("Logout failed");
    }
};
