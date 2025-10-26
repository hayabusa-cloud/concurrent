## concurrent

Concurrent is a library of lock-free and wait-free algorithms 

### Environment Requirements

Currently supports amd64 and arm64 architectures

### Basic Usage
```golang
c, p := concurrent.NewMPMCQueue[int](1024)
i := 100
err := p.Enqueue(&i)
if err != nil {
	return err
}
res, err := c.Dequeue()
if err != nil {
	return err
}
println(res)
```

### Next Step
Implement the sCQ lock-free FIFO queue

### References
- [A. Morrison and Y. Afek, "Fast concurrent queues for x86 processors," in Proc. 18th ACM SIGPLAN Symposium on Principles and Practice of Parallel Programming (PPoPP), 2013.](https://dl.acm.org/doi/10.1145/2442516.2442527)  
- [R. Nikolaev, "A scalable, portable, and memory-efficient lock-free FIFO queue," in Proc. 33rd International Symposium on Distributed Computing (DISC), 2019. LIPIcs.](https://drops.dagstuhl.de/opus/volltexte/2019/11335/pdf/LIPIcs-DISC-2019-28.pdf)  
- [N. Koval and V. Aksenov, "POSTER: Restricted memory-friendly lock-free bounded queues," in Proc. 25th ACM SIGPLAN Symposium on Principles and Practice of Parallel Programming (PPoPP), 2020, pp. 433–434.](https://nikitakoval.org/publications/ppopp20-queues.pdf)  
- [R. Nikolaev and B. Ravindran, "wCQ: A fast wait-free queue with bounded memory usage," arXiv preprint arXiv:2201.02179, Jan. 2022.](https://arxiv.org/abs/2201.02179)  
- [V. Aksenov, N. Koval, P. Kuznetsov, and A. Paramonov, "Memory bounds for concurrent bounded queues," arXiv preprint arXiv:2104.15003v5, Jan. 2024.](https://arxiv.org/abs/2104.15003)  
- [A. Denis and C. Goedefroit, "NBLFQ: A lock-free MPMC queue optimized for low contention," in Proc. 39th IEEE International Parallel and Distributed Processing Symposium (IPDPS), 2025, pp. 962–973.](https://hal.science/hal-04762608)  
- [Intel Corporation, "Combined Volume Set of Intel 64 and IA-32 Architectures Software Developer’s Manuals."](https://www.intel.com/content/www/us/en/developer/articles/technical/intel-sdm.html)  
- [Arm Limited, "Arm Architecture Reference Manual for A-profile architecture," DDI 0596, latest revision.](https://developer.arm.com/documentation/ddi0596/latest/)

### License
©2023 Hayabusa Cloud Co., Ltd.  
#5F Eclat BLDG, 3-6-2 Shibuya, Shibuya City, Tokyo 150-0002, Japan  
Released under the MIT license
