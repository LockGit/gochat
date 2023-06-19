
.PHONY:build
build: # make build TAG=1.18;make build TAG=latest,自定义版本号,构建自己机器可以运行的镜像(因为M1机器拉取提供的镜像架构不同,只能自己构建)
	cd ./docker && docker build -t lockgit/gochat:${TAG} .