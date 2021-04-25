---
title: "从C++/Go对比入门C#"
date: 2021-04-16T15:06:38+08:00
draft: true
---

自大学课程学习C++至今，已经过去了十多年。期间或课程需要，或工作需要，或自学，先后接触过C++、Perl、Objective-C、Python、Go，中间也用JS写过简陋的网站之类的。如今因为工作需要，又要加上一个C#。老实讲我不是很喜欢C++时代的语言，特性过于丰富，因此通过对比来学习新语言也是比较快速的方式了。

C#本身作为与Java同时代的竞品，身上有着C++、Java的影子。C++作为一个开放标准的高级语言，其主要目的是提供高级语言特性的同时，与C的资产无缝衔接。Java专注于跨平台、网络IO。而C#则是微软提出的一个竞品，在提供各种便利的同时，想要在Windows平台上与其自身的语言相融合。

Java和C#随着时间推移，也在各自不断引入新特性，因此能够看到一些新语言如Python/Go/Rust的同等“替代品”。

# 基础类型

C#中的基础类型与C++大同小异，如int类系强制转换也大同小异。但C#中指针类型只能在unsafe代码中使用。C#描述类型才用的主要有值类型(value type)与引用类型(reference type)。值类型在函数传参时会发生值传递，引用类型在传参时不会发生值传递。大部分基本类型都是值类型，但string除外。C#中string的机制与Java、Golang等基本一致，它指向一块只读区域，要想改写只能重新创建，或者使用专门的缓冲区组装新的字符串(`System.Text.StringBuilder`),它作为参数传递时，不会发生复制。

鉴于C/C++中默认类型转换的可能隐患，int在C++中不能用作bool类型判断，因此不能再使用类似`if(1) {}`这种语句。

Golang不允许任何不同数值类型的自动转换，因此计算前必须将两个操作数转换成指定的同类型。在C#中，这依然是一个不错的编码规范。

## 浮点数

C#的浮点数有float/double/decimal三个类型。不同于C++的不同实现定义不同，C#的浮点数长度是确定的。float为4字节，double为8字节，decimal为16字节，计算时自动向最大长度浮点类型靠拢。

三种浮点数的后缀可以分别使用f/F, d/D, m,M表示。
```csharp
double d = 3D;
d = 4d;
d = 3.934_001;

float f = 3_000.5F;
f = 5.4f;

decimal myMoney = 3_000.5m;
myMoney = 400.75M;
```

## 隐式类型推断

形如`MyPrettyClass *a = new MyPrettyClass()`这种看上去蠢蠢的代码，在旧版C++代码中随处可见。因此C++11、Golang中都加入了自动类型推断，可以通过下文内容来自动推断类型，避免写大量代码:
```c++
auto a = new MyPrettyClass();
```
如golang:
```go
a := NewMyPrettyClass()
```

C# 3.0开始，可以通过`var`关键字，通过下文来推断类型：
```csharp
var a = new MyPrettyClass();
var ia = 99999;
```

C# 9.0开始，也可以通过上文来自动推断下文的new:
```csharp
List<int> a = new(); // -> List<int> a = new List<int>();
```

`var`关键字也可以用于承接匿名类型:
```csharp
var a = new {string Name;};
```

## 数组

C++中数组不会记录长度，也不是一个单独的类型，而是首元素的指针。`sizeof`作为一个编译宏，只能返回编译时确定的长度（这里是指针长度）。因此，C++中对应数组可能更合适的是`vector<T>`。

Go中数组类型是一种值传递类型，具有确定的长度，使用场景有限。能够有动态长度的是数组切片`slice`，可以通过`len(arr)`获得存储的长度信息，作为函数参数传递时也是作为引用类型。

C#中的数组类型Array是一种**确定长度**的类型，因此声明长度必须为常量，形如：
```csharp
int[] a = new int[5];
for(int i=0; i<a.Length; i++) {
    // ...
}
```

与C++中的vector或Go中的slice接近，可以动态扩充长度的数组是`ArrayList`。

### 多维数组

```csharp
int[,] arr2D = new int[10,2];
int[,,] arr3D = new int[10,2,1];
```


### 数组的初始化和填充

C#中数组的初始化需要初始化序列与实际长度一致，不可以只初始化一部分:
```csharp
int[] a = new int[5]{1,2,3,4,5};
int[] a = {1,2,3,4,5}; // 默认创建了一个长度为5的int[]
```

数组不具备类似`memset`这种批量填充功能。


### 数组切片

在较早版本中，可以通过LINQ取出数组中的部分元素。在C# 8.0开始，可以通过[..]或[^3..]取出元素，但没有像Golang的Slice那样，可以在相同的内存区域取出数组切片。


## 析构

析构只能管理托管内存

非托管的无法在析构中自动回收


## Hello World

Golang不需要运行时，编译使用go。C#需要依赖.NET运行时，其在Windows上有完整的安装包，也可以在[Linux下安装](https://docs.microsoft.com/zh-cn/dotnet/core/install/linux-ubuntu)。

Go项目不需要特地创建项目，可以直接`go run xxx.go`。C#需要先创建项目:

`dotnet new console` 可以在当前目录下创建一个默认的控制台项目，也可以用其它模板，或者在Vs中通过wizard创建项目。这时目录下会出现一个默认的`Program.cs`:

```csharp
using System;

namespace learncs
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Hello World!");
        }
    }
}
```

## 包引入及函数入口

C#通过`using xxx`引入其它命名空间， 规则与C++基本一致，引入后具体使用就不需要再写前缀。 由于namespace可以多层嵌套，会形成和Java类似的多层引入结构。

https://docs.microsoft.com/zh-cn/dotnet/csharp/whats-new/csharp-version-history