# mdout容器化

## 镜像构建
```shell
make
```

如果不需要自己构建，可以直接使用已经构建好的镜像`lwabish/mdout`及`lwabish/mdout:chinese`

## 使用
```shell
# 将当前目录挂载到容器中，并运行mdout将当前目录中的Readme.md转换成Readme.pdf
docker run --rm -v $(pwd):/data lwabish/mdout Readme.md

# 如果markdown中有中文，需要将上述命令中的mdout镜像替换成mdout:chinese镜像
docker run --rm -v $(pwd):/data lwabish/mdout:chinese Readme.md
```

