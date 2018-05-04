require 'itamae/secrets'
secrets = Itamae::Secrets(File.join(__dir__, '../../secret'))

template '/etc/sysconfig/grass-graph' do
    owner "root"
    group "root"
    mode '0644'
    variables slack_webhook_url: secrets[:slack_webhook_url], slack_channel_name: secrets[:slack_channel_name], slack_bot_name: secrets[:slack_bot_name], google_application_credentials: secrets[:google_application_credentials]
    source "../../files/sysconfig/grass-graph.erb"
end

remote_file '/etc/systemd/system/grass-graph.service' do
    owner "root"
    group "root"
    mode '0644'
    source "../../files/systemd/grass-graph.service"
end
