package 'fuse'

remote_file '/usr/local/bin/goofys' do
    source "../../files/goofys/goofys"
    mode  '0755'
end

execute 'make directory for goofys s3 mount' do
    command <<-EOC
    mkdir -m 0755 -p /var/www/grass-graph/goofys-s3
    chown a-know:a-know /var/www/grass-graph/goofys-s3
    sudo -u a-know /usr/local/bin/goofys for-grass-graph /var/www/grass-graph/goofys-s3
    EOC
    not_if "sudo -u a-know df -h | grep goofys-s3"
end
