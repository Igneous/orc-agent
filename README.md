[![Build Status](https://drone.io/github.com/Igneous/orc-agent/status.png)](https://drone.io/github.com/Igneous/orc-agent/latest)

### What is orc-agent?

Orc agent is a generic amqp consumer that handles incoming json-serialized
messages with a user-created handler script (specified in the message's
"handler" key).

tl;dr it listens on amqp queues and forwards json payloads to handler scripts.

### Is it production ready?

Probably not for your enterprise, at least not yet. It's currently in a useable
state, and pretty easy to get up and working -- though it only ships with one
example handler (a chef handler specific to the environment I work in).

### How do I install it?

We're packaging it as a deb. Grab it from
[our drone artifact page](https://drone.io/github.com/Igneous/orc-agent/files)
and install it with dpkg.
