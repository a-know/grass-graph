require 'itamae/secrets'
secrets = Itamae::Secrets(File.join(__dir__, '../../secret'))

## install and setup td-agent
remote_file '/home/a-know/install-redhat-td-agent3.sh' do
  owner "root"
  group "root"
  mode '0755'
  source "../../files/td-agent/install-redhat-td-agent3.sh"
end

execute 'Install td-agent' do
  # See http://docs.fluentd.org/articles/install-by-rpm
  command 'sh /home/a-know/install-redhat-td-agent3.sh'
  not_if 'rpm -q td-agent'
end

service 'td-agent'

%w(
  fluent-plugin-record-reformer
  fluent-plugin-bigquery
).each do |gem|
  gem_package gem do
    gem_binary '/opt/td-agent/embedded/bin/fluent-gem'
    options '--no-rdoc --no-ri'
    notifies :restart, 'service[td-agent]'
  end
end

remote_file '/etc/td-agent/td-agent.conf' do
  owner "root"
  group "root"
  mode '0755'
  source "../../files/td-agent/td-agent.conf"
  notifies :restart, "service[td-agent]"
end

## setup files for bigquery upload
directory '/etc/td-agent/.keys'

template '/etc/td-agent/.keys/bq-credential-for-fluentd-jsonkey.json' do
  owner    'root'
  group    'root'
  mode     '0644'
  source "../../files/td-agent/bq-credential-for-fluentd-jsonkey.json.erb"
  variables(
      project_id: secrets[:bq_project_id],
      private_key_id: secrets[:bq_creds_private_key_id],
      private_key: secrets[:bq_creds_private_key],
      client_email: secrets[:bq_creds_client_email],
      client_id: secrets[:bq_creds_client_id],
      client_x509_cert_url: secrets[:bq_creds_cert_url]
  )
end

directory '/etc/td-agent/conf.d'

template '/etc/td-agent/conf.d/nginx-log.conf' do
  owner    'root'
  group    'root'
  mode     '0644'
  source "../../files/td-agent/nginx-log.conf.erb"
  variables project_id: secrets[:bq_project_id], dataset_name: node['td-agent']['nginx_log_conf']['dataset_name']
  notifies :restart, 'service[td-agent]'
end
