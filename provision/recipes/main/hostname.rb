hostname = node[:hostname]

raise 'hostname is required' unless hostname

file '/etc/hostname' do
  content "#{hostname}\n"
  notifies :run, 'execute[run hostnamectl]'
end

execute 'run hostnamectl' do
  command "hostnamectl set-hostname #{hostname}"
  action :nothing
end
