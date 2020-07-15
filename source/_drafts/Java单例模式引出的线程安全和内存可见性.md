---
title: Java单例模式引出的线程安全和内存可见性
category: Java
tags:
  - java
  - 线程安全
---

首先看一下饿汉模式的单例实现：

```java
/**
 * 饿汉模式的单例 <br>
 *
 * @author sdvdxl <杜龙少> <br>
 * @date 2020/4/21 12:58 <br>
 */
public class Singleton {

  private static final Singleton SINGLETON = new Singleton();

  private Singleton() {}

  public static Singleton getInstance() {
    return SINGLETON;
  }
}
```

现在考虑一个问题，当多线程调用 `Singleton.getInstance()` 方法的时候，是线程安全的吗？

为了更加说明问题，我们给他增加一个属性 `list`，现在完整代码如下：

```java
import java.util.ArrayList;
import java.util.List;

/**
 * 饿汉 线程安全 <br>
 *
 * @author sdvdxl <杜龙少> <br>
 * @date 2020/4/21 12:58 <br>
 */
public class Singleton {
  private static final Singleton SINGLETON = new Singleton();
  private final List<String> list;

  private Singleton() {
    this.list = new ArrayList<>();
    for (int i = 0; i < 10; i++) {
      this.list.add(String.valueOf(i));
    }
  }

  public static Singleton getInstance() {
    return SINGLETON;
  }

  public static void main(String[] args) {
    for (int i = 0; i < 10; i++) {
      int finalI = i;
      new Thread(() -> System.out.println(Singleton.getInstance().contains(String.valueOf(finalI))))
          .start();
    }
  }

  public boolean contains(String param) {
    return list.contains(param);
  }
}
```

问题是：当多线程访问 `Singleton.getInstance().contains(String.valueOf(finalI))` 的时候，是线程安全的吗？`Singleton.list` 的数据对多线程的可见性是怎样的？

首先，答案是肯定安全的，每个线程都能正确拿到list的值（里面的元素）。

假如我们把 `list` 改为 `private List<String> list;` 那结果又是如何呢？注意： list 去掉了 `final` 修饰。

没有final修饰的list变得不安全了，多线程访问不能保证list能够被安全访问（list可见性问题）。

这涉及到 final 和 安全性初始化问题。

初始化安全性将确保，对于被正确构造的对象，所有线程都能看到由构造函数为对象给哥哥 final 域设置的正确的值，而不管采用何种方式来发布对象。而且，对于可以通过正确构造对象中某个 final 域到达的任意变量，比如这里的 list 里面的元素将同样对于其他线程是可见的。

初始化安全性智能保证通过 final 域可达的值从构造过程完成时开始的可见性。对于通过非 final 域可达的值，或者在构造完成后可能改变的值，必须采用同步来确保可见性。

现在在改变一下 `list` 声明， 变为 static：

```java
public class SingletonStatic {
  private static final SingletonStatic SINGLETON = new SingletonStatic();
  private static List<String> list;

  static {
    list = new ArrayList<>();
    for (int i = 0; i < 10; i++) {
      list.add(String.valueOf(i));
    }
  }

  private SingletonStatic() {}

  public static SingletonStatic getInstance() {
    return SINGLETON;
  }
}
```

list 被 static 修饰，如果在初始化后不对list做修改动作，那么也是安全的。因为statc修饰的属性（字段）是在类加载的时候进行初始化的，只会被加载一次（正常情况下，热加载情况先不考虑），并且对多线程是可见的。

总结：final 在构造函数中正确初始化（可以理解为没有被外部线程操作），或者 被static 修饰，并且构造完成后其值（比如list的元素）不再被修改，则该对象是线程安全的；如果被 final 或者 static 修饰了，但是构造完成后，值还可以被改变则不能保证线程安全；如果没有被 final 或者 static 修饰，那么不管怎样，都是不安全的（除非该对象自身是线程安全的，比如ConcurrentHashMap）。


### 安全发布

##### 发布

“发布（PUblish）”一个对象的意思是指，使对象能够在当前作用域之外的代码中使用。在这里，比如可以通过方法 `getInstance`  将SINGLETON 发布到其他作用域（线程范围）。

发布内部状态可能会破坏封装性，并使得程序那一维持不变形。例如，如果再对象构造完成之前就发布该对象，就会破坏线程安全性。
#### 逸出

当某个不应该发布的对象被发布时，这种情况就被称为逸出（Escape)。比如懒汉模式，调用 `getInstance` 方法，可能 SINGLETON 还没有初始化完成，但是却已经可以被其他线程访问了，那么就是不安全发布，也是逸出。


```java
/**
 * 懒汉（懒加载）线程不安全 <br>
 *
 * @author sdvdxl <杜龙少> <br>
 * @date 2020/4/21 12:58 <br>
 */
public class Singleton {
  private static Singleton SINGLETON;

  private Singleton() {}

  public static Singleton getInstance() {
    if (SINGLETON == null) {
      SINGLETON = new Singleton();
    }

    return SINGLETON;
  }
}
```

#### 逸出



### 类加载

###


在singleton list例子中，list 本身并不能保证不被外部多线程修改，如果被多线程修改那么就不是线程安全的，因为 list 被发布了出去，可以在其他线程中随意使用 list的 add，addAll，remove 等能修改list元素的方法，修改后，list 的元素值极可能不会被其他线程立马更新。如果要保证list是安全的，那么有几种方式：

1. 对外提供list的副本，而不是list本身，如

```java
public List<String> getList(){
  return new ArrayList(list);
}
```

1. 使用线程安全的 list，比如 `CopyOnWriteArrayList`

```java
private final List<String> list;

private Singleton() {
    this.list = new CopyOnWriteArrayList<>();
    for (int i = 0; i < 10; i++) {
      this.list.add(String.valueOf(i));
    }
  }
```


被 volatile 修饰的变量具有内存可见性，这包含其是对象的的情况，比如：

volatile int a = 1;

vaolatile Date date;

对于a和date 来说，都是内存可见的。即使date本身没变，变的是date里面的 year参数，也是多线程可见的。

所以对于ConcurrentHashMap和CopyOnWriteArrayList等线程安全的容器来说，不仅仅能够保证容器本身的线程安全，还能保证其包含的元素的内存可见性。
