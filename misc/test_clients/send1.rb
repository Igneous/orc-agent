#!/usr/bin/env ruby
require 'json'
require 'bunny'

conn = Bunny.new 'amqp://guest:guest@172.23.25.65:5672'
conn.start

ch = conn.create_channel
tqone  = ch.queue(ARGV.first)
x  = ch.default_exchange

msg =  {
  handler:             'chef',
  run_list:            'recipe[chef-metal::nextranet]',
  override_attributes: {}
}

puts "Sending #{msg.to_json} -> #{ARGV.first}"
x.publish(msg.to_json, routing_key: tqone.name)
conn.close
