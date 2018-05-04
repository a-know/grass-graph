remote_file '/etc/systemd/system/grass-graph.service' do
    owner "root"
    group "root"
    mode '0644'
    source "../../files/systemd/grass-graph.service"
end
