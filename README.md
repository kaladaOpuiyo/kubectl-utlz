# kube-utlz

This plug-in extends the kubectl functionality by returning memory and cpu information for all pods scheduled on a Node.

## Installation

OS X
```sh
make darwin && mv ./kubectl-utlz  /usr/local/bin
```

Linux
```sh
make linux && mv ./kubectl-utlz  /usr/local/bin
```

## Usage

Basic

```sh
$ kubectl utlz ${NODE_NAME}
NAMESPACE              NAME                              CPU(cores)   MEMORY(bytes)   
dummy                  dummy-mdsj5                        80m          998Mi           
dummy                  dummy-s9hw7                        2m           27Mi            
dummy                  dummy-4x4tn                        5m           24Mi            
dummy                  dummy-95vkw                        1594m        3817Mi          
dummy                  dummy-jk99k                        113m         4002Mi          
dummy                  dummy-8qs6m                        6m           376Mi           
dummy                  dummy-wthv5                        1m           7Mi
```

Sort-By cpu/memory

```sh
$ kubectl utlz ${NODE_NAME} --sort-by=cpu
NAMESPACE              NAME                              CPU(cores)   MEMORY(bytes)   
dummy                  dummy-95vkw                        1580m        3817Mi          
dummy                  dummy-jk99k                        113m         4002Mi          
dummy                  dummy-mdsj5                        50m          998Mi           
dummy                  dummy-4x4tn                        6m           24Mi            
dummy                  dummy-8qs6m                        5m           377Mi           
dummy                  dummy-s9hw7                        1m           27Mi            
dummy                  dummy-wthv5                        1m           7Mi  
```
Views

```sh 
Resources 

$ kubectl utlz ${NODE_NAME} --view=resources 
NAMESPACE              NAME                              CPU REQUESTED(cores)   MEMORY REQUESTED(bytes)   CPU LIMITS(cores)   MEMORY LIMITS(bytes)   
dummy                  dummy-mdsj5                       1000m                  4352Mi                    0m                  4352Mi                 
dummy                  dummy-s9hw7                       10m                    0Mi                       0m                  0Mi                    
dummy                  dummy-4x4tn                       100m                   0Mi                       0m                  0Mi                    
dummy                  dummy-95vkw                       750m                   4196Mi                    5000m               8692Mi                 
dummy                  dummy-jk99k                       0m                     6144Mi                    0m                  6144Mi                 
dummy                  dummy-8qs6m                       100m                   256Mi                     256m                512Mi                  
dummy                  dummy-wthv5                       100m                   128Mi                     100m                128Mi                  
```

```sh 
Wide = resources + usage 

$ kubectl utlz ${NODE_NAME} --view=wide
NAMESPACE              NAME                              CPU(cores)   MEMORY(bytes)   CPU REQUESTED(cores)   MEMORY REQUESTED(bytes)   CPU LIMITS(cores)   MEMORY LIMITS(bytes)   
dummy                  dummy-mdsj5                        47m          998Mi           1000m                  4352Mi                    0m                  4352Mi                 
dummy                  dummy-s9hw7                        2m           27Mi            10m                    0Mi                       0m                  0Mi                    
dummy                  dummy-4x4tn                        5m           25Mi            100m                   0Mi                       0m                  0Mi                    
dummy                  dummy-95vkw                        1664m        3821Mi          750m                   4196Mi                    5000m               8692Mi                 
dummy                  dummy-jk99k                        116m         4002Mi          0m                     6144Mi                    0m                  6144Mi                 
dummy                  dummy-8qs6m                        8m           379Mi           100m                   256Mi                     256m                512Mi                  
dummy                  dummy-wthv5                        1m           7Mi             100m                   128Mi                     100m                128Mi 
```