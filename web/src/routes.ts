import page from 'page'
import type { Component } from 'svelte';

import { writable } from "svelte/store";

import LoginPage from '$modules/account/Login.svelte'
import HomePage from '$modules/Home.svelte'
import TtyPage from '$modules/system/Tty.svelte'


export const currentRoute = writable<Component>(LoginPage);
const route = (path: string, component: Component) => page(path, () => currentRoute.set(component))


route("/home", HomePage)
route("/account/login", LoginPage)
route("/tty", TtyPage)
page()
