require 'itamae/secrets'
secrets = Itamae::Secrets(File.join(__dir__, '../../secret'))

template '/home/a-know/google-service-account.json' do
    owner "root"
    group "root"
    mode '0644'
    variables project_id: secrets[:project_id], private_key_id: secrets[:private_key_id], private_key: secrets[:private_key], client_email: secrets[:client_email], client_id: secrets[:client_id], client_x509_cert_url: secrets[:client_x509_cert_url]
    source "../../files/credentials/google-service-account.json.erb"
end
