CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
tar zcvf mdout.linux.x86-64.tar.gz mdout
rm mdout

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build .
tar zcvf mdout.macOS.x86-64.tar.gz mdout
rm mdout

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build .
tar zcvf mdout_windows_x86-64.tar.gz mdout.exe
rm mdout.exe

mkdir -p release
mv mdout.linux.x86-64.tar.gz release
mv mdout.macOS.x86-64.tar.gz release
mv mdout_windows_x86-64.tar.gz release