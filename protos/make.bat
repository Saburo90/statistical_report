@echo off

D:

rem dir D:\*.* >dir.txt

cd\project\statistical_report\protos

for /r "." %%a in (*.proto) do (
    protoc -I %%~dpa --gofast_out=plugins=grpc:%%~dpa %%~nxa
)