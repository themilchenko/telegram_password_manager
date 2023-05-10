box.cfg{listen = 3301}

box.once('init', function ()
  box.schema.space.create('users')
  local users = box.space.users
  users:format({
    {name = 'chat_id', type = 'number'},
    {name = 'state', type = 'number'},
    {name = 'secret_key', type = 'string'},
    {name = 'request_service', type = 'string'}
  })
  users:create_index('primary', {
    parts = {'chat_id'}
  })

  box.schema.space.create('passwords')
  local passwords = box.space.passwords
  passwords:format({
    {name = 'chat_id', type = 'number'},
    {name = 'service_name', type = 'string'},
    {name = 'password', type = 'string'}
  })
  passwords:create_index('primary', {
    parts = {'chat_id', 'service_name'}
  })

  box.schema.user.grant('guest', 'read,write,execute', 'universe')
end)
