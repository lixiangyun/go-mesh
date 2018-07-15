export MESH_HOME=$PWD

killall controler
killall demo
killall demohttp
killall demotcp
killall mesher

cd $MESH_HOME/etcd
nohup ./etcd &

cd $MESH_HOME/bin

nohup ./controler &

nohup ./mesher -name demo -ver 1.1.1 &
nohup ./demo -n demo -v 1.0.0 &

nohup ./mesher -name demohttp -ver 1.0.0 &
nohup ./demohttp -n demohttp -v 1.0.0 -p 127.0.0.1:8001 &

nohup ./mesher -name demohttp -ver 1.0.1 &
nohup ./demohttp -n demohttp -v 1.0.1 -p 127.0.0.1:9001 &

nohup ./mesher -name demotcp -ver 1.0.0 &
nohup ./demotcp -n demotcp -v 1.0.0 -p 127.0.0.1:10001 &
