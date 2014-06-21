#!/usr/bin/env ruby
require 'json'
require 'bunny'

conn = Bunny.new 'amqp://orc-agent:famc@10.10.10.10:5672'
conn.start

ch = conn.create_channel
q  = ch.queue('a4:ba:db:fd:9c:52')
x  = ch.default_exchange

obj = {
  handler:             'chef',
  run_list:            'recipe[webapp_stack]',
  override_attributes: {
    aws: {
      owner: 'bwolfe',
      id:    42
    }
  }
}

x.publish(obj.to_json, routing_key: q.name)

sleep 1.0

conn.close
