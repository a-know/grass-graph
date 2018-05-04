directory '/var/www/grass-graph/app' do
    owner 'a-know'
    group 'a-know'
    mode  '0755'
    action :create
end

directory '/var/www/grass-graph/log' do
    owner 'a-know'
    group 'a-know'
    mode  '0755'
    action :create
end
