execute 'create swapfile' do
    user 'root'
    command <<-EOC
      dd if=/dev/zero of=/swap.img bs=1M count=2048 &&
      chmod 600 /swap.img
      mkswap /swap.img
      swapon /swap.img
    EOC
    only_if "test ! -f /swap.img -a `cat /proc/swaps | wc -l` -eq 1"
end
