groups = {}

user node['user']['uid'] do
  gid node['user']['gid'] unless node['user']['uid'] == node['user']['gid']
  password node['user']['password']
  home node['user']['home']
  shell node['user']['shell']
end

(node['user']['groups'] || []).each do |name|
  groups[name] ||= []
  groups[name] << node['user']['uid']
end

directory node['user']['home'] do
  owner node['user']['uid']
  group node['user']['gid']
  mode '0755'
  action :create
end

[ "#{node['user']['home']}/.ssh" ].each do |dir|
  directory dir do
    owner node['user']['uid']
    group node['user']['gid']
    mode '0755'
    action :create
  end
end

file "#{node['user']['home']}/.ssh/authorized_keys" do
  content node['user']['authorized_keys'].join("\n")
  owner node['user']['uid']
  group node['user']['gid']
  mode '0600'
  action :create
end

groups.each do |name, members|
  group name do
    members members
    action :manage
  end
end

execute "add to sudoers" do
  user "root"
  command "echo '#{node['user']['uid']} ALL=NOPASSWD: ALL' >> /etc/sudoers"
  not_if "grep #{node['user']['uid']} /etc/sudoers"
end
