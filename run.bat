set MESH_HOME=%CD%

cd %MESH_HOME%\etcd
start .\etcd

cd %MESH_HOME%\bin

start .\controler.exe

start .\mesher.exe -name demo -ver 1.1.1
start .\demo.exe -n demo -v 1.0.0

start .\mesher.exe -name demohttp -ver 1.0.0
start .\demohttp.exe -n demohttp -v 1.0.0 -p 127.0.0.1:8001

start .\mesher.exe -name demohttp -ver 1.0.1
start .\demohttp.exe -n demohttp -v 1.0.1 -p 127.0.0.1:9001

start .\mesher.exe -name demotcp -ver 1.0.0
start .\demotcp.exe -n demotcp -v 1.0.0 -p 127.0.0.1:10001

pause