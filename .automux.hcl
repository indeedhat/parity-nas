session = "parity-nas"
# config = "./tmux.conf"
# attach_existing = false

window "Editor" {
    exec = "vim"
    focus = true
}

window "Shell" {
    split {
        dir = "web/"
        exec = "npm run dev"
    }
    split {
        vertical = true
        exec = "air"
    }
}
