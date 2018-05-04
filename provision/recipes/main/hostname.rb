hostname = node[:hostname]

raise 'hostname is required' unless hostname

execute 'update hostname' do
  hostname_regexp = Regexp.escape(hostname)
  command <<-EOC
    sed -i -e 's/\\(HOSTNAME=\\).*/\\1#{hostname_regexp}/' /etc/sysconfig/network
    hostname #{hostname}
  EOC
  not_if "hostname | grep #{hostname}"
end

remote_file '/etc/profile.d/network.sh' do
  source "../../files/hostname/network.sh"
  action :create
end
