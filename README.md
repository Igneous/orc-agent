### Dafuq is orc-agent?

Orc agent is a generic amqp consumer that handles incoming json-serialized
messages with a user-created script (based off of the message's handler value).

tl;dr it listens on amqp queues and forwards json payloads to handler scripts.
(and if you couldn't tell, this is completely alpha and incomplete)
(don't use it, yet)
