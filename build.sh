export MESH_HOME=$PWD

if [ -e bin ]; then
    rm -rf bin
fi

mkdir bin

cd $MESH_HOME/controler
go build .
mv controler $MESH_HOME/bin
cp config.json $MESH_HOME/bin

cd $MESH_HOME/mesher
go build .
mv mesher $MESH_HOME/bin

cd $MESH_HOME/example/demo
go build .
mv demo $MESH_HOME/bin

cd $MESH_HOME/example/demohttp
go build .
mv demohttp $MESH_HOME/bin

cd $MESH_HOME/example/demotcp
go build .
mv demotcp $MESH_HOME/bin
