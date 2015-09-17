# go-numato

A simple library for controlling Numato (numato.com) relay and GPIO
devices. Also provides a test double for testing code that relies on
communication with a Numato controller.

## Limitations (TODO)

* Currently only fully supports operations with Relays (my common use case).
* Doesn't auto-detect device capability.
* Device test can only be run with a relay board that has at least 4 relays. I
like the sound, deal with it.

## License

Under MIT license, see LICENSE file.
