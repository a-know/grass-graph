template '/etc/selinux/config' do
  source "../../files/selinux/selinux.config.centos.erb"
  variables policy: 'permissive'
  notifies :run, 'execute[setenforce_permissive]'
end
  
execute 'setenforce_permissive' do
  action :nothing
  command <<-EOC
if [ $(getenforce) != "Disabled" ]; then
  setenforce permissive
fi
EOC
end
  