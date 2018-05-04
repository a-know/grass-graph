service 'crond'
service 'rsyslog'

zone = 'Asia/Tokyo'

file '/etc/localtime' do
  action :delete
  not_if 'readlink -f /etc/localtime | grep Tokyo'
end

link '/etc/localtime' do
  to "/usr/share/zoneinfo/#{zone}"
  notifies :restart, 'service[crond]', :immediately
  notifies :restart, 'service[rsyslog]'
end
