start ..\..\mesher\mesher.exe -name demohttp -ver 1.0.0
start demohttp.exe -n demohttp -v 1.0.0 -p 127.0.0.1:8001
start ..\..\mesher\mesher.exe -name demohttp -ver 1.0.1
demohttp.exe -n demohttp -v 1.0.1 -p 127.0.0.1:9001
