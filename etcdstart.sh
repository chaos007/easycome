../../etcd/etcd.exe --name infra0 --initial-advertise-peer-urls http://192.168.10.89:2380 \
--listen-peer-urls http://192.168.10.89:2380 \
--listen-client-urls http://192.168.10.89:2379,http://127.0.0.1:2379 \
--advertise-client-urls http://192.168.10.89:2379 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster infra0=http://192.168.10.89:2380 \
--initial-cluster-state new