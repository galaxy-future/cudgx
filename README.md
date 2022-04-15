![画板备份 9](https://user-images.githubusercontent.com/94337797/148376608-8cc8efe6-dd60-44db-9eb3-e7e93dc329c3.png)
-------


[![Go Report Card](https://goreportcard.com/badge/github.com/galaxy-future/BridgX)](https://goreportcard.com/report/github.com/galaxy-future/BridgX) &nbsp;
[![CodeFactor](https://www.codefactor.io/repository/github/galaxy-future/bridgx/badge)](https://www.codefactor.io/repository/github/galaxy-future/bridgx)


Language
----

English | [中文](https://github.com/galaxy-future/cudgx/blob/master/docs/CH-README.md)

Introduction
-----
CudgX is an AIOps intelligent operation & maintenance engine for the cloud-native era launched by Galaxy Future, it indexes and digitally measures various services through multi-dimensional big-data collection and machine learning training analysis. Based on the deployed training model, the real-time measurement of service quality could achieve the automated and intelligent real-time dynamic scheduling of computing power, storage, network, and other basic resources. 



Main Features:<br>
1.	Support service logging record and aggregate computing；<br>
2.	Support automatic pressure measurement 
3.	Support real-time capacity evaluation 
4.	Support automatic, intelligent computing-force elastic schedule
5.	Provide agency access to mesh agent SDK
6.	Support automatic measurement and Auto Scaling of Web service metric
7.	With the open platform of metrics definition, partners could configure customized metrics based on the platform.


Contact Us
----

[Weibo](https://weibo.com/galaxyfuture) | [Zhihu](https://www.zhihu.com/org/xing-yi-wei-lai) | [Bilibili](https://space.bilibili.com/2057006251)
| [WeChat Official Account](https://github.com/galaxy-future/comandx/blob/main/docs/resource/wechat_official_account.md)
| [WeCom Communication Group](https://github.com/galaxy-future/comandx/blob/main/docs/resource/wechat.md)


System Architecture
--------

* System Architecture Diagram

<img width="839" alt="image" src="https://user-images.githubusercontent.com/94337797/153759267-c96e8a8f-9771-4f53-b20a-22cc7a45917b.png">


* System Module Description

  * cudgx-api：
    * predict-rules： Responsible for maintaining Scaling rules.

    * redundancy-keeper： Responsible for maintaining the redundancy of service clusters.

  * cudgx-gateway： Responsible for collecting dotting data from `metrics-go` and distributing it to Kafka.
  * cudgx-consumer： Responsible for consuming Kafka data, and storing it in clickhouse.
  * metrics-go： CudgX dots SDK.

System Diagram Legends
--------
- View Service Monitoring Legends: Redundancy Trends, QPS, Machine Amount 
- View record of service scaling 
- Example:

<img width="713" alt="image" src="https://user-images.githubusercontent.com/94337797/153759376-b2fa4699-c439-4d3e-b513-029971126512.png">


Installation and Deployment 
----
1.	Configuration Requirements
To ensure the stable running of the system, it is recommended that the system has 2 cores and 4 GB memory storage; CudgX has been installed and tested in Linux and MacOS.  

2.	Environment Dependency 
Before installing CudgX, please install:
-	BridgX: Please install BridgX according to [BridgX installation Guide](https://github.com/galaxy-future/bridgx/blob/master/README.md). The intranet deployment environment must be able to communicate with the cloud vendor’s VPC
-	SchedulX: Please install SchedulX according to the [SchedulX installation Guide](https://github.com/galaxy-future/schedulx/blob/master/README.md). The intranet deployment environment must be able to communicate with the cloud vendor’s VPC
-	ComandX: If front-end operations are required, please install ComandX according to the [ComandX installation Guide] (https://github.com/galaxy-future/comandx/blob/main/README.md).


3.	Installment Steps
* (1) Download the source code:
   - Back-end engineering:
  > `git clone https://github.com/galaxy-future/cudgx.git`

* (2)MacOS(intel chip only) system installation and deployment 

    - Back-end deployment, running in CudgX directory

      > `make docker-run-mac`

* (3)Linux system installation and deployment 

    - 1）For users:

        - Back-end deployment, running in CudgX directory 

          > `make docker-run-linux`

    - 2）For developers 

        - Since the project will download the required basic images, it is recommended to place the downloaded source code in a directory that is larger than 10 GB
        - Back-end deployment
            - CudgX relies on mysql & kafka & clickhouse artifact,
                - If use the built-in mysql & etcd & clickhouse, access the CudgX root directory, and use the following command:

                  > docker-compose up -d    //start up CudgX <br>
                  > docker-compose down    //stop CudgX  <br>
                - If already have external mysql & etcd & clickhouse service, go to  `cd conf` and change the IP and port configuration information for the corresponding artifacts of `api.json` `gateway.json` `consumer.json`, then go to the root directory of CudgX, use the following command:
                  > docker-compose up -d api    //start up api <br>
                  > docker-compose up -d gateway //tart up gateway <br>
                  > docker-compose up -d consumer //tart up consumer <br>
                  > docker-compose down     //stop CudgX

* (4)Front-end page
        
  - If front-end operations are required, please install [ComandX](https://github.com/galaxy-future/comandx/blob/main/README.md)
    -	After the system runs, enter 'http://127.0.0.1' to view the management console page. The initial username is “root” and the password is “123456”.


* (5)metrics-go
    - 1）To complete uploading the data, the target application should dot based on CudgX-SDK [metrics-go](https://github.com/galaxy-future/metrics-go/blob/master/README.md), we have provided a sample application `cudgx-sample-pi` for testing: The `cudgx-sample-pi` application has been buried based on metrics-go, and the docker image `galaxyfuture/cudgx-sample-pi` has been pushed to docker hub. 
    
    
    - 2）Completing service deployment through SchedulX, see the process of creating SchedulX services. 

      - Key Deployment:

        - ComandX Page - Service Deployment - Create Scaling Process – Image Deployment – Services Startup Command:
            > docker run -d -e CUDGX_SERVICE_NAME=test.cudgx.gf -e CUDGX_CLUSTER_NAME=default -e CUDGX_GATEWAY_URL=http://127.0.0.1:8080 -p 80:8090  
        
            Environment Variables Parameter Description：
 
          - CUDGX_SERVICE_NAME： Sverice Name
          - CUDGX_CLUSTER_NAME： Cluster Name
          - CUDGX_GATEWAY_URL： CudgX-gateway Service IP address

        - ComandX Page-Service Deployment - Create Scaling Process – Traffic Access – Configure SLB ID: Connect to Ali cloud SLB.
  
    - 3）Rules of Scaling Configuration
      
      - Key Configuration
        - ComandX Page – Service Deployment – Scaling Rules 


    - 4）Two Methods are recommended for pressuring {SLB IP}:{SLB PORT}/pi Interface
   
      - a：During the test, we provided `cudgx-sample-benchmark` application to simulate access traffic. Docker images `galaxyfuture/cudgx-sample-benchmark` has been pushed to docker hub
        
        > docker run -d --name cudgx_sample_benchmark --network host galaxyfuture/cudgx-sample-benchmark --gf.cudgx.sample.benchmark.sever-address={SLB IP}:{SLB PORT}/pi
            
      - b：By using opensource interface testing tools, apply pressure to {SLB IP}:{SLB PORT}/pi 

    - 5）ComandX Page- Service List – Cluster Monitoring 
      - View service monitoring legends: Redundancy Trends, QPS, Machine Amount 
      -	View service scaling records    





Code of Conduct
------
[Contributor Convention](https://github.com/galaxy-future/cudgx/blob/master/CODE_OF_CONDUCT.md)

Authorization
-----

CudgX uses [Elastic License 2.0](https://github.com/galaxy-future/cudgx/blob/master/LICENSE) Agreement for Authorization.

Contact us
-----
If you want more information about the service, scan the following QR code to contact us:
![image](https://user-images.githubusercontent.com/102009012/163563581-45874bce-2f79-46ee-b309-2757ddd17c46.png)

