## concurrent

Concurrent is a library of lock-free and wait-free algorithms 

### Environment Requirements

Currently only supporting amd64 architecture

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
* [Combined Volume Set of Intel® 64 and IA-32 Architectures Software Developer’s Manuals](https://www.intel.com/content/www/us/en/developer/articles/technical/intel-sdm.html)
* [Ruslan Nikolaev and Binoy Ravindran. 2022. wCQ: A Fast Wait-Free Queue with Bounded Memory Usage](https://arxiv.org/pdf/2201.02179.pdf)
* [Nikita Koval and Vitaly Aksenov. 2020. POSTER: Restricted Memory-Friendly Lock-Free Bounded Queues](https://nikitakoval.org/publications/ppopp20-queues.pdf)
* [Ruslan Nikolaev. 2019. A Scalable, Portable, and Memory-Efficient Lock-Free FIFO Queue](https://drops.dagstuhl.de/opus/volltexte/2019/11335/pdf/LIPIcs-DISC-2019-28.pdf)
* [Adam Morrison and Yehuda Afek. 2013. Fast Concurrent Queues for x86 Processors](https://dl.acm.org/doi/10.1145/2442516.2442527)

### License
©2023 Hayabusa Cloud Co., Ltd.  
#5F Eclat BLDG, 3-6-2 Shibuya, Shibuya City, Tokyo 150-0002, Japan  
Released under the MIT license
