set MESH_HOME=%CD%

rmdir /q/s bin
mkdir bin

cd %MESH_HOME%\controler
go build .
move .\controler.exe %MESH_HOME%\bin
copy .\config.json %MESH_HOME%\bin

cd %MESH_HOME%\mesher
go build .
move .\mesher.exe %MESH_HOME%\bin

cd %MESH_HOME%\example\demo
go build .
move .\demo.exe %MESH_HOME%\bin

cd %MESH_HOME%\example\demohttp
go build .
move .\demohttp.exe %MESH_HOME%\bin

cd %MESH_HOME%\example\demotcp
go build .
move .\demotcp.exe %MESH_HOME%\bin
