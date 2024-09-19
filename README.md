# 随心游小程序需求文档
## 1 引言

### 1.1 编写目的

&emsp;&emsp;本文档的目的是详细地介绍随心游小程序其Alpha版本所包含的需求，以便客户能够确认产品的确切需求以及开发人员能够根据需求设计编码。

### 1.2 背景
&emsp;&emsp;该项目为旅游工具性质的小程序，接口文档正在github创建，前端正在开发中。

## 2 任务概述

### 2.1 项目概述

#### 2.1.1 项目来源及背景
&emsp;&emsp;随心游微信小程序是一个可根据用户性格、心情、出发地、目的地等因素推荐短途旅行地点，规划旅行行程的旅行工具。用户可使用小程序在汕头规划自己的短途旅行的路线

#### 2.1.2 项目目标
&emsp;&emsp;我们希望该小程序功能尽量精简实用，页面简洁，用户能简单通过页面文字使用所需的功能，高效准确地为用户定制服务。


### 2.1.3 系统功能概述
1. 关于用户

(1)注册与登录：小程序自动通过微信openID进行注册与登录。

(2)设置信息：用户可以查看或修改头像、个性签名、性别、电话号码、性格等基础信息。

(3)查看信息：用户可以查看如基本信息、文章总量以及列表、点赞、旅行路径收藏。

(4)搜索用户：用户可以搜索其它用户。

(5)用户留言：创建、查看、删除留言。

(6)用户好友：添加、通过、拒绝、删除、拉黑。

（7）历史足迹： 用户可以查看走过的规划旅程。


2.关于用户组

(1)管理用户组：用户可以创建、更新、查看、删除用户组信息。

(2)申请加入：用户可以申请加入某个用户组，组长也可以通过、拒绝申请。

(3)黑名单：用户组可以拉黑名单，黑名单人员无法加入用户组。

(4)点赞/踩、收藏：用户可以点赞/踩某个用户组。

(5)标签：用户可以给用户组创造、查看、删除标签。

(5)搜索用户组：用户可以搜索其它用户组。

(6)用户组列表：用户可以查看热度排名、点赞排名、收藏排名等。

(7)标准用户组：用户可以创建标准用户组。

3.关于社区

(1)文章：用户可以对文章进行增删查改、点赞、收藏、查看游览数量并搜索、查看文章排行。

(2)公告：发布、订阅、查看公告。

4.关于聊天

(1)私信：用户可以与其它用户私信。

## 3 页面设计



### 3.1 全局
**导航栏** ：底部导航栏包含 **首页** **定制** **我的** 页面跳转

**搜索框** ：**首页**和**定制**页面有顶部搜索框，不可被隐藏

### 3.2 首页

**搜索框** ：首页顶部有顶部搜索框，可搜索文章信息。搜索文章后可跳转到结果**文章列表**

**公告栏** ：搜索框下有滚动播放栏，点击播放窗口图片可跳转推荐文章或者发布的公告。
                可用来推荐文章

**文章推荐列表**：文章推荐，点击后可跳转到具体文章内容。
向下滑寻找文章可以隐藏**滚动播放栏**

**文章列表**：**搜索框**搜索结果列表

**文章页面**： 文章页面包括文章内容，图片，以及**收藏**、**点赞**、**留言**功能

### 3.3 定制

**地图窗口**： 用户定位（申请用户位置信息的权限），规划路径

**搜索框**： 底部有搜索栏，可搜索目的地，以及设置起始点（默认起始点为用户位置）、心情、路程、交通工具等信息要素

&emsp;&emsp;或者选择以**起始点**推荐目的地

&emsp;&emsp;无论是**搜索目的地**还是选择以**起始点**推荐目的地，都进行路径规划，在地图上显示路径规划的界面。生成的路径可导出为图片进行分享

**目的地推荐**： 搜索框搜索目的地后，用图标显示目的地

**导航**： 用户确定目的地和路径之后，为用户进行导航

### 3.4 我的

**个性化**：展示用户头像、账号、昵称（如无特殊设置一律为微信名）、格言

**个人信息**： 展示用户手机号、用户头像、账号、昵称（如无特殊设置一律为微信名）、格言，**可修改用户信息**

**历史足迹**：展示用户历史足迹

**收藏**： 包括用户的路径收藏和文章收藏

**系统公告**

## 3 功能需求

### 3.1 功能描述

**1.用户相关**

(1)用户注册

功能描述：小程序通过微信服务器提供的API获取用户openID，通过openID为用户注册

&emsp;&emsp;使用接口：

(2)用户登录
功能描述：用户登录小程序，为前端页面返回tokenID凭证，提供使用用户注销、创建、修改文章的权限

&emsp;&emsp;使用接口：

(3)修改用户信息
功能描述：

&emsp;&emsp;使用接口：

(4)查看用户信息
功能描述：显示当前用户性别、手机号、IP地址等信息

&emsp;&emsp;使用接口：

（5）用户收藏
功能描述：用户可收藏文章或者足迹

&emsp;&emsp;使用接口：

（6）查看收藏
功能描述：用户可查看收藏的文章或者足迹

&emsp;&emsp;使用接口：

（7）查看历史足迹
功能描述：用户可以查看走过的规划旅程。

（8）公告
功能描述：管理员账户可发布、查看、更新、删除公告公告，公告可能需要拥有一定的引导、跳转能力。



**2.社区相关**

(1)文章

功能描述：用户可以发布文章、查看文章、点赞收藏文章、文章标签、文章分类等。

&emsp;&emsp;使用接口：

（2）文章评论

&emsp;&emsp;使用接口：

（3）文章转发
功能描述：用户可以转发文章

&emsp;&emsp;使用接口：

**3.路径规划相关**

（1）定位

功能描述：定位用户当前位置

（2）目的地推荐

功能描述：输入用户起点和终点、旅程范围和时间等数据为用户推荐一个或多个目的地和途径地点，并且为用户规划最佳路径

（3）路径图片生成

功能描述：用户可将系统生成的路径生成图片

（3）路径分享

功能描述：用户可以分享历史足迹

# 4 其他需求

### 4.1 验收标准

交付周期为**两周一交付**，交付后由测试人员依据用户故事进行功能测试，并提出改进方案以及可能存在的bug。Alpha版本不包括exam系统、competition系统以及file系统。**迭代次数不得超过6次**。

### 4.2 资源建设

在交付周期进行时，我们需要制作出至少一份面向大众的页面引导以及各类资源。

### 4.3 对外联通

我们会尝试链接大众点评、美团、交通出行等功能

### 4.4 推荐系统

推荐系统要求实现输入一组用户的文章历史浏览信息、用户地址，推荐相关的文章；输入出发地、目的地、性格、心情等一组数据，推荐相关地点和规划最优路径

