tb
====

tb is tiny http balance which I made for fun.

# Usage

You can run it in docker. To configure it, You have to create a file in ``/etc/tb.yml`` in the container.
So, an example of configuration is placed in ``./conf.yml``


# Goals

- [x] configuration file
- [ ] weighted round-robin/least connections
- [ ] optimize searching of alive backends

---

Inspired by: [kasvith](https://github.com/kasvith/simplelb/)
