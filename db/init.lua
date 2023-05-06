local box = require('box')
box.cfg{listen = 3301}

box.schema.space.create('user_states')
local user_states = box.space.user_states
user_states:format({
  {name = 'user_id', type = 'number'},
  {name = 'state', type = 'number'}
})
user_states:create_index('primary', {
  type = 'hash',
  parts = {'user_id'}
})

box.schema.space.create('passwords')
local passwords = box.space.passwords
passwords:format({
  {name = 'user_id', type = 'number'},
  {name = 'service_name', type = 'string'},
  {name = 'password', type = 'string'}
})
passwords:create_index('primary', {
  type = 'hash',
  parts = {'user_id', 'service_name'}
})
