import page from 'page'
import type { Component } from 'svelte';

import { writable } from "svelte/store";

import LoginPage from '$modules/account/Login.svelte'
import HomePage from '$modules/Home.svelte'
import TtyPage from '$modules/system/Tty.svelte'
import { guardAdminUser, guardGuest, guardUser } from './lib/auth';


type Middleware = (ctx, next) => void


export const currentRoute = writable<Component>(LoginPage);

const route = (path: string, component: Component, ...middleware: Middleware[]) => {
    page(path, ...(middleware || []), () => currentRoute.set(component))
}


route("/home", HomePage, guardUser)
route("/account/login", LoginPage, guardGuest)
route("/tty", TtyPage, guardAdminUser)

page()
