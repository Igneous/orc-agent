#!/usr/bin/env ruby
require 'json'
require 'bunny'

conn = Bunny.new 'amqp://orc-agent:ownme@mq.your.com:5672'
conn.start

ch = conn.create_channel
q  = ch.queue(ARGV.first)
x  = ch.default_exchange

msg =  {
  handler:             'chef',
  run_list:            'recipe[your_webapp:redeploy]',
  override_attributes: { branch: 'ticket-3902',
                         treeish: '93a7de' }
}

puts "Sending #{msg.to_json} -> #{ARGV.first}"
x.publish(msg.to_json, routing_key: q.name)
conn.close
