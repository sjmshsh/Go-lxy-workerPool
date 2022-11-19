# Go-lxy-workerPool
轻量级线程池

虽然Goroutine的开销非常廉价，但也不是免费的。Go1.4版本开始采用了连续栈的方法，也就是每一个Goroutine的执行栈都是一块儿连续的内存，如果空间不足，运行时会分配一个更大的连续内存空间作为整个Goroutine的执行栈，将原栈内容拷贝到新分配的空间中来。

虽然整个方案可以避免Go1.3采用的分段栈会导致的**hot split问题**，但连续栈的原理也决定了，一旦Goroutine的执行栈发生了grow，那么即便这个Goroutine不再需要这么大的栈空间，这个Goroutine的栈空间也不会被Shrink(收缩)了，这些空间可能会处于长时间闲置的状态，直到Goroutine退出。

另一方面，Go运行时进行Goroutine调度的处理器消耗，也会随之增加，称为阻碍Go应用性能提升的重要因素。

> 这种问题常见的解决方法就是：使用Gouroutine池，把M个计算任务调度到N个Goroutine上，而不是为每个计算任务分配一个独享的Goroutine，从而提高计算资源的利用率

### workerpool的实现原理

- pool的创建与销毁
- pool中worker(Goroutine)的管理
- task的提交与调度

![workerpool架构图.excalidraw](C:\Users\lenovo\Desktop\图片\板书\workerpool架构图.excalidraw.png)

capacity是pool的一个属性，代表整个pool中worker的最大容量。我们使用一个带缓冲的channel：active，作为worker的**计数器**，这种channel的使用模式就是我们之前讲过的**计数信号量**。

当active channel可写的时候，我们就创建一个worker，用于处理用户通过Schedule函数提交的待处理的请求。当active channel满了的时候，pool就会停止worker的创建，直到某个worker因故退出，active channel又空出一个位置的时候，pool才会创建新的worker填补那个空位。
