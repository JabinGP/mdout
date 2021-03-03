# mdout 安装指南

## 0. 特殊国情

由于大陆的特殊国情，访问 github 速度很慢甚至会失败，以下的 https 链接都来自 github ，如果访问失败可以尝试使用我提供的代理下载。

代理的使用方式为，在原有的完整链接之前加上 `https://ghproxy.cfjabin.workers.dev/` 。

例如 github 原链接为 `https://github.com/JabinGP/mdout/releases/download/v0.6/mdout.macOS.x86-64.tar.gz` ，加上代理前缀后则是 `https://ghproxy.cfjabin.workers.dev/https://github.com/JabinGP/mdout/releases/download/v0.6/mdout.macOS.x86-64.tar.gz` 。

该代理支持访问仓库源码文件，也支持对仓库的克隆，release文件的下载。

但需要注意的是，该代理会将一些加速分散到其他免费的cdn服务上，例如 `jsDeliv` ，因此无法保证通过该代理获取的文件与 github 原链接的文件是一致的，因为cdn存在更新延迟。

我无法保证该代理的可用性（即能成功访问并下载），也无法保证该代理的一致性（即通过代理下载的文件与原链接文件一致），通常情况下使用该代理不会出现什么问题，但请以 github 原链接为准。

## 1. 安装方式说明

mdout 的安装分为三部分：

1. 获取 mdout 可执行文件
2. 配置 mdout 命令
3. 初始化 mdout 配置文件、主题

如果你是 macOS 或者 linux ，这三个部分可以被简化成一个脚本：

```cmd
bash -c "$(wget https://raw.githubusercontent.com/FisherWY/Shell/master/mdout/install_mdout.sh -O -)"
```

### 2. 获取 mdout 可执行文件与命令配置

mdout 可执行文件有两种方式获取：

1. 源码编译
2. 直接从release界面获取可执行文件

#### 2.1. 源码编译并配置命令

> 从源码编译要求 golang 1.13 以上版本。

```bash
# 克隆源码
git clone https://github.com/JabinGP/mdout.git
# 进入项目
cd mdout
# 编译
go build .
# 执行成功在项目目录中生成了可执行文件 mdout
```

#### 2.2. 从release下载对应平台可执行文件

mdout 基于 golang ，得益于 golang 交叉编译的特点，可以提前为各平台（windows、linux、macOS）打包可执行文件release出来供大家下载，这样就不需要要求使用者具备 golang 开发环境了。

进入 [release界面](https://github.com/JabinGP/mdout/releases) 找到最新的 `release` ，打开 `Asserts` 既可以看到类似如下的文件列表：

- mdout.linux.x86-64.tar.gz 7.82 MB
- mdout.macOS.x86-64.tar.gz 7.72 MB
- mdout_windows_x86-64.tar.gz 7.65 MB
- Source code (zip)
- Source code (tar.gz)

下载 mdoutxxxxx 解压后即可获得对应平台的可执行文件。

##### 2.2.1 windows获取可执行文件并配置的步骤

1. 选择 `mdout_windows_x86-64.tar.gz` 左键点击下载，下载完成后解压即可获得 `mdout.exe`

2. 确定  `mdout.exe` 所在文件夹，例如我的电脑中的 `D:\mdout`
    ![1](./markdown/1.jpg)  

3. 设置环境变量，右键我的电脑 -> 选择属性 -> 左边的高级系统设置
    ![2](./markdown/2.jpg)  

4. 选择高级 -> 点击环境变量
    ![3](./markdown/3.jpg)  

5. 找到下半部分的系统变量，双击`Path`行
    ![4](./markdown/4.jpg)  

6. 在弹出来的窗口选择新建
    ![5](./markdown/5.jpg)  

7. 填入`D:\mdout`，然后一定要连续点完三个确定
    ![6](./markdown/6.jpg)  

8. 检验是否成功打开cmd，或者powershell，或者你有gitbash都ok（推荐使用命令行的windows用户都装一个gitbash），输入`mdout`，看到如下输出就是成功了
    ![7](./markdown/7.jpg)

##### 2.2.2 linux获取可执行文件并配置命令的步骤

- 选择 `mdout.linux.x86-64.tar.gz` 右键复制链接

1. 通过命令行下载、解压

    ```bash
    wget {$DownloadLink}
    tar -xvzf mdout.linux.x86-64.tar.gz
    ```

2. 放入系统可执行文件目录

    ```bash
    sudo mv mdout /usr/local/bin
    ```

3. 检验是否成功

    ```bash
    # 看到版本号输出就成功
    mdout --version
    ```

##### 2.2.3 macOS获取可执行文件并配置命令的步骤

- 选择 `mdout.macOS.x86-64.tar.gz` 右键复制链接

1. 通过命令行下载、解压

    ```bash
    wget {$DownloadLink}
    tar -xvzf mdout.macOS.x86-64.tar.gz
    ```

2. 放入系统可执行文件目录

    ```bash
    mv mdout /usr/local/bin
    ```

3. 检验是否成功

    ```bash
    # 看到版本号输出就成功
    mdout --version
    ```

### 3. 初始化配置文件和主题包

mdout 的配置文件和主题有两种方式初始化：

1. 通过 mdout 内置命令获取网络资源到本地
2. 手动创建配置文件和主题

#### 3.1. 配置文件夹说明

mdout 在运行时会自动在用户的家目录下创建一个名为 mdout 的配置文件夹，一般而言结构应该如下所示：

```text
mdout
|-- conf.toml
|-- log
|   `-- xxxx-x.log
`-- theme
    |-- github
        |-- css
        |-- index.html
        `-- lib
```

不同平台的家目录不同，最终导致不同平台的 mdout 配置文件夹路径不同，以下表格以用户为 `jabin` 举例：

|平台|配置文件夹路径|
|-|-|
|windows|C:/Users/jabin/mdout|
|linux|/home/jabin/mdout|
|macOS|/Users/jabin/mdout|

#### 3.2. 通过 mdout 内置命令获取网络资源到本地

##### 3.2.1. 通过 mdout 内置命令获取配置文件 `conf.toml`

```bash
mdout install config
```

该命令会默认获取代码中预设的 github 文件地址，如果出现网络问题，或者预设的文件不对，可以通过 `-u` 参数指定一个配置文件链接下载。

```bash
mdout install config -u {$ConfigLink}
```

|版本|配置文件链接|
|-|-|
|v0.6|[github链接](https://raw.githubusercontent.com/JabinGP/mdout/v0.6/asserts/config/conf.toml)|

##### 3.2.2. 通过 mdout 内置命令获取配套主题包

下载主题包需要指定两个参数：

1. 第一个参数为 `-u` 接上主题包下载链接
2. 第二个参数 `-n` 指定该主题包下载后命名为什么，推荐名为 `mdout` ，因为不修改配置文件的情况下，mdout默认会使用名为 `gihub` 的主题包。

```cmd
mdout install theme -u {$ThemeDownloadLink} -n {${ThemeName}}
```

|主题|仓库地址|主题包下载链接|
|-|-|-|
|仿github主题v0.6-0.1|[JabinGP/mdout-theme-github](https://github.com/JabinGP/mdout-theme-github)|[github链接](https://github.com/JabinGP/mdout-theme-github/releases/download/v0.6-0.1/mdout-theme-github-v0.6-0.1.zip)|

#### 3.3. 手动创建配置文件和主题

参考 3.1. 中的不同平台配置文件夹路径以及文件位置：

1. 创建 mdout 文件夹。
2. 参考项目源码中 `assets/config/conf.toml` 创建自己的 `conf.toml`。
3. 创建 theme 文件夹，利用 git 将主题包克隆到该文件夹中。

    ```bash
    # 例如将 JabinGP/mdout-theme-github 克隆并保存为主题名 github
    git clone https://github.com/JabinGP/mdout-theme-github github
    ```
