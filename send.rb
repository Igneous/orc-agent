#!/usr/bin/env ruby
require 'json'
require 'bunny'

conn = Bunny.new 'amqp://orc-agent:famc@10.10.10.10:5672'
conn.start

ch = conn.create_channel
tqone  = ch.queue('testqueue1')
tqtwo  = ch.queue('testqueue2')
x  = ch.default_exchange

def msg(dst, count)
  {
    handler:             'chef',
    run_list:            'recipe[webapp_stack]',
    override_attributes: {
      aws: {
        owner: 'bwolfe',
        id:    42,
        time: Time.now,
        dst: dst.to_s,
        count: count
      }
    }
  }
end

1000.times do |count|
  sample = [ :heads, :tails ].sample

  if sample == :heads
    "#{count} -> testqueue1"
    x.publish(msg(sample, count).to_json, routing_key: tqone.name)
  else
    "#{count} -> testqueue2"
    x.publish(msg(sample, count).to_json, routing_key: tqtwo.name)
  end
end

sleep 1.0

conn.close
