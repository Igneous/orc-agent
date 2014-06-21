#!/usr/bin/env ruby
require 'bunny'

conn = Bunny.new 'amqp://orc-agent:famc@10.10.10.10:5672'
conn.start

ch = conn.create_channel
q  = ch.queue('a4:ba:db:fd:9c:52')
x  = ch.default_exchange

1000.times do |time|
  x.publish("(#{time + 1}) Hello! The time is #{Time.now}", :routing_key => q.name)
end

sleep 1.0

conn.close
