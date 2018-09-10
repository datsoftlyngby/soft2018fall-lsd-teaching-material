We have a small application, a telephone book application. It consists of a [webserver](webserver/telbook_server.go) written in Go, which finds all `Person` records stored in a MongoDB [database](dbserver/db_setup.js) and returns an HTML page with those records.


# An example with Vagrant

See the [Vagrantfile](Vagrantfile) and start it with `vagrant up`. It will create two VMs. one with the [webserver](main.go) and the other one with the [database](db_setup.js).

# An example with Docker

## Building the DB Server
The following is written from my perspective, i.e. user `helgecph`.

```bash
$ cd dbserver/
$ docker build -t helgecph/dbserver .
```

## Building the Webserver

```bash
$ cd webserver/
$ docker build -t helgecph/webserver .
```

Now, check that both images are locally available.

```bash
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
webserver           latest              5520fac0a523        24 seconds ago      718MB
dbserver            latest              f5567a451a4a        2 minutes ago       359MB
```

## Starting the Application Manually



```bash
$ mkdir $(pwd)/datadb
$ docker run -d -p 27017:27017 --name dbserver helgecph/dbserver
$ docker run -it -d --rm --name webserver --link dbserver -p 8080:8080 helgecph/webserver
```

Eventhough deprecated, on can `--link` the containers via the bridge network together.

```bash
$ docker ps -a
CONTAINER ID        IMAGE                COMMAND                  CREATED             STATUS              PORTS                      NAMES
0282fc8b2c41        helgecph/webserver   "/bin/sh -c ./telb..."   11 seconds ago      Up 10 seconds       0.0.0.0:8080->8080/tcp     webserver
06b85924f444        helgecph/dbserver    "docker-entrypoint..."   6 minutes ago       Up 6 minutes        0.0.0.0:27017->27017/tcp   dbserver
```

Properly done, from now on on links containers via a shared network.

```bash
$ docker network create example-network
$ docker network ls
NETWORK ID          NAME                            DRIVER              SCOPE
d5a8f5d3b2c2        bridge                          bridge              local
9c9d24069da7        example-network                 bridge              local
bd11ae20c3ac        host                            host                local
51892d4cc44a        none                            null                local
$ docker run -d -p 27017:27017 --name dbserver --network=example-network helgecph/dbserver
$ docker run -it -d --rm --name webserver --network=example-network -p 8080:8080 helgecph/webserver
```


### Testing the Application

```bash
$ docker run --rm --network=example-network appropriate/curl:latest curl http://webserver:8080
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0<!DOCTYPE HTML>
<html>
    <head>
        <title>The Møllers</title>
    </head>
    <body>
        <h1>Telephone Book</h1>
        <hr>
        <table style="width:50%">
          <tr>
            <th>Index</th>
            <th>Name</th>
            <th>Phone</th>
            <th>Address</th>
            <th>City</th>
          </tr>

          <tr>
            <td>0</td>
            <td>Møller</td>
            <td>&#43;45 20 86 46 44</td>
            <td>Herningvej 8</td>
            <td>4800 Nykøbing F</td>
          </tr>

          <tr>
            <td>1</td>
            <td>A Egelund-Møller</td>
            <td>&#43;45 54 94 41 81</td>
            <td>Rønnebærparken 1 0011</td>
            <td>4983 Dannemare</td>
          </tr>

          <tr>
            <td>2</td>
            <td>A K Møller</td>
            <td>&#43;45 75 50 75 14</td>
            <td>Bregnerødvej 75, st. 0002</td>
            <td>3460 Birkerød</td>
          </tr>

          <tr>
            <td>3</td>
            <td>A Møller</td>
            <td>&#43;45 97 95 20 01</td>
            <td>Dalstræde 11 Heltborg</td>
            <td>7760 Hurup Thy</td>
          </tr>

        </table>
        <p></p>
        Data taken from <a href="https://www.krak.dk/person/resultat/møller">Krak.dk</a>
    </body>
</html>
100  1366  100  1366    0     0   443k      0 --:--:-- --:--:-- --:--:--  666k
```


## Stopping the Application Manually


```bash
$ docker stop dbserver
$ docker stop webserver
```

```bash
$ docker rm webserver
$ docker rm dbserver
```

## Starting the Application with Docker Compose


```yml
version: '3'
services:
  dbserver:
    image: helgecph/dbserver
    ports:
      - "27017:27017"
    networks:
      - outside

  webserver:
    image: helgecph/webserver
    ports:
      - "8080:8080"
    networks:
        - outside

  clidownload:
    image: appropriate/curl
    networks:
      - outside
    entrypoint: sh -c  "sleep 5 && curl http://webserver:8080"

networks:
  outside:
    external:
      name: example-network
```


```bash
$ docker-compose up
Creating 03-containers_and_vms_clidownload_1 ... done
Creating 03-containers_and_vms_webserver_1   ... done
Creating 03-containers_and_vms_dbserver_1    ... done
Attaching to 03-containers_and_vms_clidownload_1, 03-containers_and_vms_dbserver_1, 03-containers_and_vms_webserver_1
dbserver_1     | 2018-09-09T13:15:35.219+0000 I CONTROL  [main] Automatically disabling TLS 1.0, to force-enable TLS 1.0 specify --sslDisabledProtocols 'none'
dbserver_1     | about to fork child process, waiting until server is ready for connections.
dbserver_1     | forked process: 26
dbserver_1     | 2018-09-09T13:15:35.220+0000 I CONTROL  [main] ***** SERVER RESTARTED *****
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] MongoDB starting : pid=26 port=27017 dbpath=/data/db 64-bit host=f3c46139bd3d
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] db version v4.0.1
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] git version: 54f1582fc6eb01de4d4c42f26fc133e623f065fb
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] OpenSSL version: OpenSSL 1.0.2g  1 Mar 2016
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] allocator: tcmalloc
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] modules: none
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] build environment:
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten]     distmod: ubuntu1604
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten]     distarch: x86_64
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten]     target_arch: x86_64
dbserver_1     | 2018-09-09T13:15:35.227+0000 I CONTROL  [initandlisten] options: { net: { bindIp: "127.0.0.1", port: 27017, ssl: { mode: "disabled" } }, processManagement: { fork: true, pidFilePath: "/tmp/docker-entrypoint-temp-mongod.pid" }, systemLog: { destination: "file", logAppend: true, path: "/proc/1/fd/1" } }
dbserver_1     | 2018-09-09T13:15:35.227+0000 I STORAGE  [initandlisten]
dbserver_1     | 2018-09-09T13:15:35.227+0000 I STORAGE  [initandlisten] ** WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine
dbserver_1     | 2018-09-09T13:15:35.227+0000 I STORAGE  [initandlisten] **          See http://dochub.mongodb.org/core/prodnotes-filesystem
dbserver_1     | 2018-09-09T13:15:35.227+0000 I STORAGE  [initandlisten] wiredtiger_open config: create,cache_size=3479M,session_max=20000,eviction=(threads_min=4,threads_max=4),config_base=false,statistics=(fast),log=(enabled=true,archive=true,path=journal,compressor=snappy),file_manager=(close_idle_time=100000),statistics_log=(wait=0),verbose=(recovery_progress),
dbserver_1     | 2018-09-09T13:15:35.775+0000 I STORAGE  [initandlisten] WiredTiger message [1536498935:775355][26:0x7fb4e774ca00], txn-recover: Set global recovery timestamp: 0
dbserver_1     | 2018-09-09T13:15:35.780+0000 I RECOVERY [initandlisten] WiredTiger recoveryTimestamp. Ts: Timestamp(0, 0)
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten] ** WARNING: Access control is not enabled for the database.
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten] **          Read and write access to data and configuration is unrestricted.
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten] ** WARNING: /sys/kernel/mm/transparent_hugepage/enabled is 'always'.
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten] **        We suggest setting it to 'never'
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten] ** WARNING: /sys/kernel/mm/transparent_hugepage/defrag is 'always'.
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten] **        We suggest setting it to 'never'
dbserver_1     | 2018-09-09T13:15:35.788+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:35.789+0000 I STORAGE  [initandlisten] createCollection: admin.system.version with provided UUID: ad940693-6c94-4667-9ce8-8a2f254de608
dbserver_1     | 2018-09-09T13:15:35.795+0000 I COMMAND  [initandlisten] setting featureCompatibilityVersion to 4.0
dbserver_1     | 2018-09-09T13:15:35.797+0000 I STORAGE  [initandlisten] createCollection: local.startup_log with generated UUID: 860095e6-0b3c-4ef4-ab81-7c125a515170
dbserver_1     | 2018-09-09T13:15:35.805+0000 I FTDC     [initandlisten] Initializing full-time diagnostic data capture with directory '/data/db/diagnostic.data'
dbserver_1     | 2018-09-09T13:15:35.806+0000 I NETWORK  [initandlisten] waiting for connections on port 27017
dbserver_1     | child process started successfully, parent exiting
dbserver_1     | 2018-09-09T13:15:35.807+0000 I STORAGE  [LogicalSessionCacheRefresh] createCollection: config.system.sessions with generated UUID: 0e19c266-255c-4c6d-bd46-6187125419e6
dbserver_1     | 2018-09-09T13:15:35.821+0000 I INDEX    [LogicalSessionCacheRefresh] build index on: config.system.sessions properties: { v: 2, key: { lastUse: 1 }, name: "lsidTTLIndex", ns: "config.system.sessions", expireAfterSeconds: 1800 }
dbserver_1     | 2018-09-09T13:15:35.821+0000 I INDEX    [LogicalSessionCacheRefresh] 	 building index using bulk method; build may temporarily use up to 500 megabytes of RAM
dbserver_1     | 2018-09-09T13:15:35.821+0000 I INDEX    [LogicalSessionCacheRefresh] build index done.  scanned 0 total records. 0 secs
dbserver_1     | 2018-09-09T13:15:35.853+0000 I NETWORK  [listener] connection accepted from 127.0.0.1:37162 #1 (1 connection now open)
dbserver_1     | 2018-09-09T13:15:35.854+0000 I NETWORK  [conn1] received client metadata from 127.0.0.1:37162 conn1: { application: { name: "MongoDB Shell" }, driver: { name: "MongoDB Internal Client", version: "4.0.1" }, os: { type: "Linux", name: "Ubuntu", architecture: "x86_64", version: "16.04" } }
dbserver_1     | 2018-09-09T13:15:35.856+0000 I NETWORK  [conn1] end connection 127.0.0.1:37162 (0 connections now open)
dbserver_1     |
dbserver_1     | /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/db_setup.js
dbserver_1     | 2018-09-09T13:15:35.895+0000 I NETWORK  [listener] connection accepted from 127.0.0.1:37164 #2 (1 connection now open)
dbserver_1     | 2018-09-09T13:15:35.896+0000 I NETWORK  [conn2] received client metadata from 127.0.0.1:37164 conn2: { application: { name: "MongoDB Shell" }, driver: { name: "MongoDB Internal Client", version: "4.0.1" }, os: { type: "Linux", name: "Ubuntu", architecture: "x86_64", version: "16.04" } }
dbserver_1     | 2018-09-09T13:15:35.899+0000 I STORAGE  [conn2] createCollection: test.people with generated UUID: e6f9ea55-18df-41ed-8630-ade8d42336ec
dbserver_1     | 2018-09-09T13:15:35.914+0000 I NETWORK  [conn2] end connection 127.0.0.1:37164 (0 connections now open)
dbserver_1     |
dbserver_1     |
dbserver_1     | 2018-09-09T13:15:35.930+0000 I CONTROL  [main] Automatically disabling TLS 1.0, to force-enable TLS 1.0 specify --sslDisabledProtocols 'none'
dbserver_1     | killing process with pid: 26
dbserver_1     | 2018-09-09T13:15:35.934+0000 I CONTROL  [signalProcessingThread] got signal 15 (Terminated), will terminate after current cmd ends
dbserver_1     | 2018-09-09T13:15:35.934+0000 I NETWORK  [signalProcessingThread] shutdown: going to close listening sockets...
dbserver_1     | 2018-09-09T13:15:35.934+0000 I NETWORK  [signalProcessingThread] removing socket file: /tmp/mongodb-27017.sock
dbserver_1     | 2018-09-09T13:15:35.934+0000 I CONTROL  [signalProcessingThread] Shutting down free monitoring
dbserver_1     | 2018-09-09T13:15:35.935+0000 I FTDC     [signalProcessingThread] Shutting down full-time diagnostic data capture
dbserver_1     | 2018-09-09T13:15:35.935+0000 I STORAGE  [signalProcessingThread] WiredTigerKVEngine shutting down
dbserver_1     | 2018-09-09T13:15:36.004+0000 I STORAGE  [signalProcessingThread] shutdown: removing fs lock...
dbserver_1     | 2018-09-09T13:15:36.005+0000 I CONTROL  [signalProcessingThread] now exiting
dbserver_1     | 2018-09-09T13:15:36.005+0000 I CONTROL  [signalProcessingThread] shutting down with code:0
dbserver_1     |
dbserver_1     | MongoDB init process complete; ready for start up.
dbserver_1     |
dbserver_1     | 2018-09-09T13:15:36.952+0000 I CONTROL  [main] Automatically disabling TLS 1.0, to force-enable TLS 1.0 specify --sslDisabledProtocols 'none'
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] MongoDB starting : pid=1 port=27017 dbpath=/data/db 64-bit host=f3c46139bd3d
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] db version v4.0.1
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] git version: 54f1582fc6eb01de4d4c42f26fc133e623f065fb
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] OpenSSL version: OpenSSL 1.0.2g  1 Mar 2016
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] allocator: tcmalloc
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] modules: none
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] build environment:
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten]     distmod: ubuntu1604
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten]     distarch: x86_64
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten]     target_arch: x86_64
dbserver_1     | 2018-09-09T13:15:36.957+0000 I CONTROL  [initandlisten] options: { net: { bindIpAll: true } }
dbserver_1     | 2018-09-09T13:15:36.957+0000 I STORAGE  [initandlisten] Detected data files in /data/db created by the 'wiredTiger' storage engine, so setting the active storage engine to 'wiredTiger'.
dbserver_1     | 2018-09-09T13:15:36.957+0000 I STORAGE  [initandlisten]
dbserver_1     | 2018-09-09T13:15:36.957+0000 I STORAGE  [initandlisten] ** WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine
dbserver_1     | 2018-09-09T13:15:36.957+0000 I STORAGE  [initandlisten] **          See http://dochub.mongodb.org/core/prodnotes-filesystem
dbserver_1     | 2018-09-09T13:15:36.957+0000 I STORAGE  [initandlisten] wiredtiger_open config: create,cache_size=3479M,session_max=20000,eviction=(threads_min=4,threads_max=4),config_base=false,statistics=(fast),log=(enabled=true,archive=true,path=journal,compressor=snappy),file_manager=(close_idle_time=100000),statistics_log=(wait=0),verbose=(recovery_progress),
dbserver_1     | 2018-09-09T13:15:37.657+0000 I STORAGE  [initandlisten] WiredTiger message [1536498937:656996][1:0x7f5d68d0ea00], txn-recover: Main recovery loop: starting at 1/26112
dbserver_1     | 2018-09-09T13:15:37.745+0000 I STORAGE  [initandlisten] WiredTiger message [1536498937:745930][1:0x7f5d68d0ea00], txn-recover: Recovering log 1 through 2
dbserver_1     | 2018-09-09T13:15:37.803+0000 I STORAGE  [initandlisten] WiredTiger message [1536498937:803760][1:0x7f5d68d0ea00], txn-recover: Recovering log 2 through 2
dbserver_1     | 2018-09-09T13:15:37.849+0000 I STORAGE  [initandlisten] WiredTiger message [1536498937:849679][1:0x7f5d68d0ea00], txn-recover: Set global recovery timestamp: 0
dbserver_1     | 2018-09-09T13:15:37.858+0000 I RECOVERY [initandlisten] WiredTiger recoveryTimestamp. Ts: Timestamp(0, 0)
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten] ** WARNING: Access control is not enabled for the database.
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten] **          Read and write access to data and configuration is unrestricted.
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten] ** WARNING: /sys/kernel/mm/transparent_hugepage/enabled is 'always'.
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten] **        We suggest setting it to 'never'
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten] ** WARNING: /sys/kernel/mm/transparent_hugepage/defrag is 'always'.
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten] **        We suggest setting it to 'never'
dbserver_1     | 2018-09-09T13:15:37.867+0000 I CONTROL  [initandlisten]
dbserver_1     | 2018-09-09T13:15:37.879+0000 I FTDC     [initandlisten] Initializing full-time diagnostic data capture with directory '/data/db/diagnostic.data'
dbserver_1     | 2018-09-09T13:15:37.880+0000 I NETWORK  [initandlisten] waiting for connections on port 27017
clidownload_1  |   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
clidownload_1  |                                  Dload  Upload   Total   Spent    Left  Speed
dbserver_1     | 2018-09-09T13:15:40.106+0000 I NETWORK  [listener] connection accepted from 172.19.0.3:48954 #1 (1 connection now open)
100  136<!DOCTYPE HTML> 0     0      0      0 --:--:-- --:--:-- --:--:--     0
clidownload_1  | <html>
clidownload_1  |     <head>
clidownload_1  |         <title>The M?llers</title>
clidownload_1  |     </head>
clidownload_1  |     <body>
clidownload_1  |         <h1>Telephone Book</h1>
clidownload_1  |         <hr>
clidownload_1  |         <table style="width:50%">
clidownload_1  |           <tr>
clidownload_1  |             <th>Index</th>
clidownload_1  |             <th>Name</th>
clidownload_1  |             <th>Phone</th>
clidownload_1  |             <th>Address</th>
clidownload_1  |             <th>City</th>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>0</td>
clidownload_1  |             <td>M?ller</td>
clidownload_1  |             <td>&#43;45 20 86 46 44</td>
clidownload_1  |             <td>Herningvej 8</td>
clidownload_1  |             <td>4800 Nyk?bing F</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>1</td>
clidownload_1  |             <td>A Egelund-M?ller</td>
clidownload_1  |             <td>&#43;45 54 94 41 81</td>
clidownload_1  |             <td>R?nneb?rparken 1 0011</td>
clidownload_1  |             <td>4983 Dannemare</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>2</td>
clidownload_1  |             <td>A K M?ller</td>
clidownload_1  |             <td>&#43;45 75 50 75 14</td>
clidownload_1  |             <td>Bregner?dvej 75, st. 0002</td>
clidownload_1  |             <td>3460 Birker?d</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>3</td>
clidownload_1  |             <td>A M?ller</td>
clidownload_1  |             <td>&#43;45 97 95 20 01</td>
clidownload_1  |             <td>Dalstr?de 11 Heltborg</td>
clidownload_1  |             <td>7760 Hurup Thy</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |         </table>
clidownload_1  |         <p></p>
clidownload_1  |         Data taken from <a href="https://www.krak.dk/person/resultat/m?ller">Krak.dk</a>
clidownload_1  |     </body>
clidownload_1  | </html>
dbserver_1     | 2018-09-09T13:15:40.108+0000 I NETWORK  [conn1] end connection 172.19.0.3:48954 (0 connections now open)
clidownload_1  | 6  100  1366    0     0   190k      0 --:--:-- --:--:-- --:--:--  190k
03-containers_and_vms_clidownload_1 exited with code 0
dbserver_1     | 2018-09-09T13:15:51.692+0000 I NETWORK  [listener] connection accepted from 172.19.0.3:48956 #2 (1 connection now open)
dbserver_1     | 2018-09-09T13:15:51.693+0000 I NETWORK  [conn2] end connection 172.19.0.3:48956 (0 connections now open)
dbserver_1     | 2018-09-09T13:15:51.708+0000 I NETWORK  [listener] connection accepted from 172.19.0.3:48958 #3 (1 connection now open)
dbserver_1     | 2018-09-09T13:15:51.709+0000 I NETWORK  [conn3] end connection 172.19.0.3:48958 (0 connections now open)
```

### Cleaning up

```bash
$ docker ps -a
CONTAINER ID        IMAGE                COMMAND                  CREATED             STATUS                       PORTS               NAMES
01a0a11d00d3        appropriate/curl     "sh -c 'sleep 5 &&..."   9 minutes ago       Exited (0) 9 minutes ago                         03containersandvms_clidownload_1
ef4617bdc0d8        helgecph/webserver   "/bin/sh -c ./telb..."   9 minutes ago       Exited (137) 6 seconds ago                       03containersandvms_webserver_1
113c782030c4        helgecph/dbserver    "docker-entrypoint..."   9 minutes ago       Exited (0) 5 seconds ago                         03containersandvms_dbserver_1
```

```bash
$ docker-compose rm
``




## Before Cleaning-up Containers

```bash
$ docker images
REPOSITORY           TAG                 IMAGE ID            CREATED             SIZE
helgecph/dbserver    latest              7672dc76725d        11 seconds ago      359MB
helgecph/webserver   latest              7e20874fe656        30 seconds ago      718MB
mongo                latest              b39de1d79a53        2 weeks ago         359MB
appropriate/curl     latest              f73fee23ac74        3 weeks ago         5.35MB
golang               jessie              6ce094895555        4 weeks ago         699MB
```


