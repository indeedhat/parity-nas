# proxy
- [ ] sub domains
- [x] enable/disable rewrite features per proxy

# logger
- [ ] rotate log files
- [ ] live viewer for logs (probably driven by server side events)
  - technically this exists but its ugly, it needs some propper formatting

# docker
- [ ] provide a web ui to manage docker compose based applications

# samba
- [ ] provide a web ui to manage samba shares

# tty
- [x] persist the terminal session between page loads

# plugins
- [x] provide a way to extend api functionalyty with go plugins
- [ ] provide a way to extend ui routes via a fs provided by plugin
  - [ ] this may need to recompile the ui if i stick with using svelte
    - i might swap to htmx as i don't think svele actually gives me any benefit for this project other than the component libs

# health check
- [ ] run a health check on startup to ensure all required dependencies are present on the machine

# moon shot
- [ ] create a custom device driver to provide parity on non striped drives

# cron
- [ ] setup a scheduled task system

# cli
- [ ] provide a cli interface for common tasts

# cleanup
- [ ] reduce external dependencies as much as possible
- [ ] make sure everything is providing propper logging

# docs
- [ ] write some
