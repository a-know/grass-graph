remote_file '/var/tmp/setup_repo.sh' do
  owner    'root'
  group    'root'
  mode     '0755'
  source   "../../files/mackerel/setup_repo.sh"
end

execute 'setup mackerel yum repo' do
  command 'sh /var/tmp/setup_repo.sh'
  not_if 'test -e /etc/yum.repos.d/mackerel.repo'
end

package 'mackerel-agent'
package 'mackerel-agent-plugins'
package 'mackerel-check-plugins'

service 'mackerel-agent'

directory '/etc/mackerel-agent/conf.d'

require 'itamae/secrets'
secrets = Itamae::Secrets(File.join(__dir__, '../../secret'))


template '/etc/mackerel-agent/mackerel-agent.conf' do
  owner 'root'
  group 'root'
  mode '0644'
  source "../../files/mackerel/mackerel-agent.conf.erb"
  variables apikey: secrets[:mackerel_api_key], service_role: node['mackerel']['service_role']
  notifies :restart, 'service[mackerel-agent]'
end

remote_file '/usr/lib/systemd/system/mackerel-agent.service' do
  owner 'root'
  group 'root'
  mode '0644'
  source "../../files/mackerel/mackerel-agent.service"
  notifies :run, 'execute[run daemon-reload]'
end

execute 'run daemon-reload' do
  command "sudo systemctl daemon-reload"
  action :nothing
end

remote_file '/etc/mackerel-agent/conf.d/check-plugins.conf' do
  owner 'root'
  group 'root'
  mode '0644'
  source "../../files/mackerel/check-plugins.conf"
  notifies :restart, 'service[mackerel-agent]'
end

remote_file '/etc/mackerel-agent/conf.d/metric-plugins.conf' do
  owner 'root'
  group 'root'
  mode '0644'
  source "../../files/mackerel/metric-plugins.conf"
  notifies :restart, 'service[mackerel-agent]'
end

# for kazeburo script plugin
package 'perl'

remote_file '/etc/mackerel-agent/periodic-checker' do
  owner 'root'
  group 'root'
  mode '0755'
  source "../../files/mackerel/periodic-checker"
end
  