# Parity NAS
This will eventually become a NAS management web UI solution

It is predicated on the existance of a custom linux multi device driver that i am currently in the process of learning
how to write and as such regardless the completion level of the interface in this repository it is largely usless.

## Motivation
Learning. I have not really done much with linux device drivers in the past and wanted to dig a little deeper into
how to create them, i run a home nas/media server and this seemed like a good opportunity to see what i could acomplish.

## Wishlist
Im currently debating the pros and cons of creating handling this project as a full os or stack that can be installed 
on any linux os and a seperate distro.  
Im proceeding with the expectation of it being a stack but that could well change in the future.

- [ ] Custom multi device driver
    - I want to be able to mount multiple drives with a standard unstriped filesystem but also provide parity
- [ ] Mobile/desktop friendly UI
    - Primerally focused on desktop management but it would be nice if the UI supported doing simple operations
      in a mobile friendly manner.
- [.] System monitor
- [ ] Samba shares/user management
- [ ] Docker compose support for running custom apps
    - It would be nice to have an integration to access the docker web ui's using a pari-nas ui  
      for example plex at nas.local/_/plex rather than something like 192.168.1.104:32400
- [.] common config interface for all parinas functionality
