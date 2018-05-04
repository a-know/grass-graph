package 'logrotate'

remote_file '/etc/logrotate.d/grass-graph' do
  owner "root"
  group "root"
  mode '0644'
  source "../../files/logrotate/grass-graph"
end
