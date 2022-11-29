<p align="center">
<img src="ui/src/assets/logos/logo-horizontal.svg" alt="drawing" width="500"/>
    </p>
    
# podlox: Skinny Kubernetes Traffic Viewer<br /> (including local traffic)

<p align="center">
    <a href="https://github.com/mogenius/podloxx/blob/main/LICENSE">
        <img alt="GitHub License" src="https://img.shields.io/github/license/mogenius/podloxx?logo=GitHub&style=flat-square">
    </a>
    <a href="https://github.com/mogenius/podloxx/releases/latest">
        <img alt="GitHub Latest Release" src="https://img.shields.io/github/v/release/mogenius/podloxx?logo=GitHub&style=flat-square">
    </a>
    <a href="https://github.com/mogenius/podloxx/releases">
      <img alt="GitHub all releases" src="https://img.shields.io/github/downloads/mogenius/podloxx/total">
    </a>
    <a href="https://github.com/mogenius/podloxx">
      <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/mogenius/podloxx">
    </a>
    <a href="https://discord.gg/WSxnFHr4qm">
      <img alt="Discord" src="https://img.shields.io/discord/932962925788930088?logo=mogenius">
    </a>
</p>

<p align="center">
  <img src="assets/screenshot1.png" alt="drawing" width="800"/>
</p>
<br />
<br />

# TOC
- [What is podloxx?](#what-is-podloxx)
- [Installation](#installation)
- [How does it work? What does it do?](#how-does-it-work-what-does-it-do)
- [Configuration](#configuration)
- [API](#api)
- [Roadmap](#roadmap)
- [Known Issues](#known-issues)
- [Thanks](#thanks)

# What is podloxx?
podloxx is a easy to use cluster daemonset and frontend application to monitor the ongoing traffic of your cluster in real-time. The monitoring happens in real-time and the display happens every few seconds (for perfromance reasons).
The most important goal was to tell local and external traffic appart. This allows you to identify high traffic applications and see how it communicates with the outside world. This even works with slim containers :-)

# Installation
Just download it and run it. Don't forget to set the right cluster using kubectx or whatever tool you prefer.

## Mac
```
curl -Lo podloxx https://github.com/mogenius/podloxx/releases/download/v1.0.11/podloxx-1.0.11-darwin-arm64 && chmod 755 podloxx

podloxx start
```
## Linux
```
curl -Lo podloxx https://github.com/mogenius/podloxx/releases/download/v1.0.11/podloxx-1.0.11-linux-amd64 && chmod 755 podloxx

podloxx start
```
## Windows
```
curl -Lo podloxx https://github.com/mogenius/podloxx/releases/download/v1.0.11/podloxx-1.0.11-win-amd64 && chmod 755 podloxx

podloxx start
```

⚠️ ⚠️ ⚠️ IMPORTANT: be sure to select the right context before running podloxx ⚠️⚠️⚠️

```
podloxx start
```

# How does it work? What does it do?
Podloxx will run a series of tasks in order to run within your cluster:
1. Create a podloxx namespace (to separate it from other workloads)
2. Setup RBAC (for proper access control)
3. Setup a memory-only redis (all DaemonSets will drop theire data here)
4. Create a DaemonSet (scrape data from all nodes)
5. Create redis-service to make the redis accessible via port forward (service makes it easier for podloxx)
6. Setup port forward to redis-service
7. Start a webservice locally to expose a web app (which gathers the data from redis)
8. Start a browser to open the web app

In other words: The daemonset will inspect all packets of the node (using special deployment capabilities). The data will be captured, summarized and send to redis (using certain thresholds). The local web app will gather the data from the redis periodically and display the data inside a fancy UI in your web browser.

As soon as you close the cli app (CTRL + C) the application will be removed from your cluster and the UI will stop receiving updates. When you restart it, it will resume gathering data without storing a state (meaning your start from 0). You can also run following command to clean up everything:
```
podloxx clean
```

# Configuration
To support multiple CNI configuration we provided a parameter to setup. "--interface-prefix"
```
podloxx start --interface-prefix azv|veth|cali
```

| Provider      | CNI         | Prefix    | Tested|
| ------------- |:----------- |:---------:| -----:|
| Azure         | Azure CNI   |       azv |    👍 |
| Azure         | -           |      veth |    👍 |
| Azure         | Calico      |      cali |    👍 |
| AWS           | CNI         |       - |      ❓ |
| AWS           | -           |       - |      ❓ |
| AWS           | Calico      |       - |      ❓ |
| Google Cloud  | CNI         |       - |      ❓ |
| Google Cloud  | -           |       - |      ❓ |
| Google Cloud  | Calico      |       - |      ❓ |

If you have tested a different configuration: let us know what works :-)

# API
You can use following API endpoints to use the raw data:
```
http://127.0.0.1:1337/traffic/overview
http://127.0.0.1:1337/traffic/total
http://127.0.0.1:1337/traffic/flow
```


# Roadmap
XXX 

# Known Issues
- Sometimes the pordforward does not get established and we do not recognize it. We are trying to fix this.
- Depending on your cloud provider your INTERFACE_PREFIX might be different. If you encounter a different interface please report its name to us so we can improve our list.

# Thanks 
We took great inspiration (and some lines of code) from [Mizu](https://github.com/up9inc/mizu).</br>
Awesome work from the folks at [UP9](https://up9.com/).</br>
Notice: The project has been renamed to Kubeshark and moved to https://github.com/kubeshark/kubeshark.</br>
