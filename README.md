<p align="center">
  <img src="ui/src/assets/logos/logo-horizontal.svg" alt="drawing" width="500"/>
</p>

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

> :warning: **This is the first public Podloxx release. If you experience issues we are very happy about your feedback or contribution.**

# Table of contents
- [What is Podloxx?](#what-is-podloxx)
- [Installation](#installation)
- [How does Podloxx work?](#how-does-podloxx-work)
- [Configuration](#configuration)
- [Known Issues](#known-issues)
- [Credits](#credits)

# What is podloxx?
Podloxx is a light-weight Kubernetes traffic monitoring tool that can be deployed as daemonset in every cluster. Local and external traffic is monitored in real time and can be viewed through a simple web interface. This allows identification of high traffic applications, understanding container relations and optimizing Kubernetes setups. Even works with slim containers 🙃

# Installation
Just download it and run it. Don't forget to set the right cluster using kubectx or whatever tool you prefer.

## Mac
```

podloxx_link=$(curl -s https://api.github.com/repos/mogenius/podloxx/releases/latest | grep browser_download_url | cut -d '"' -f 4 | grep darwin  )

curl -s -L -o podloxx ${podloxx_link}

chmod 755 podloxx

./podloxx start
```

## Linux
```

podloxx_link=$(curl -s https://api.github.com/repos/mogenius/podloxx/releases/latest | grep browser_download_url | cut -d '"' -f 4 | grep linux )

curl -s -L -o podloxx ${podloxx_link}

chmod 755 podloxx

./podloxx start
```

## Windows
```

curl.exe -LO "https://github.com/mogenius/podloxx/releases/download/v1.0.4/podloxx-1.0.4-windows-amd64"
podloxx-1.0.4-windows-amd64 start

```

⚠️ ⚠️ ⚠️ IMPORTANT: be sure to select the right context before running podloxx ⚠️⚠️⚠️

```
./podloxx start
```

# How does Podloxx work?
Podloxx will run a series of tasks in order to run within your cluster. Here's what happens in detail once you launch Podloxx:
1. A Podloxx namespace is created to isolate it from other workloads.
2. Set up RBAC for proper access control.
3. Start a memory-only redis. All DaemonSets will drop their data here.
4. Create a DaemonSet to scrape data from all nodes.
5. Launch a redis service to make the redis accessible via port forwarding.
6. Set up port forwarding for the redis service.
7. Start a web service locally to expose the Podloxx web application (which gathers the data from redis).
8. Launch the web application in a browser.

In other words: The DaemonSet will inspect all packages of the node (using special deployment capabilities). The data will be captured, summarized and sent to the redis (using certain thresholds). The local web app will gather the data from the redis periodically and display the data inside the web application.

As soon as you close the cli app (CTRL + C) the application will be removed from your cluster and the UI will stop receiving updates. When you restart it, it will resume gathering data without storing a state (meaning you start from 0).

To completely remove Podloxx from your cluster run:
```
./podloxx clean
```

# TESTED WITH
We already checked multiple CNI configurations.

| Provider      | CNI         | Prefix    | K8S    | Tested|
| ------------- |:----------- |:---------:|:---------:| -----:|
| Azure         | Azure CNI   |       azv | 1.24.X, 1.23.x, 1.22.x |    👍 |
| Azure         | -           |      veth | 1.24.X, 1.23.x, 1.22.x |    👍 |
| Azure         | Calico      |      cali | 1.24.X, 1.23.x, 1.22.x |    👍 |
| DigitalOcean  | Cillium     |       lxc | 1.24.X, 1.23.X         |    👍 |
| AWS           | CNI         |       eni | 1.24.X, 1.23.x, 1.22.x |    👍 |
| AWS           | -           |       - |         - |      ❓ |
| AWS           | Calico      |       - |         - |      ❓ |
| Google Cloud  | CNI         |       - |         - |      ❓ |

If you have tested additional configurations: Let us know what works :-)
💥: 1.25.X is not yet supported (at least we saw a problem with Digital Ocean) because the CONFIG_CGROUP_PIDS flag is disabled by default.

# API
You can use following API endpoints to access the raw data:
```
http://127.0.0.1:1337/traffic/overview
http://127.0.0.1:1337/traffic/total
http://127.0.0.1:1337/traffic/flow
```

# Known Issues
- Sometimes port forwarding doesn't get established and Podloxx doesn't recognize it. Please just hit CTRL + C to recover from this state.

# Credits
We took great inspiration (and some lines of code) from [Mizu](https://github.com/up9inc/mizu).</br>
Awesome work from the folks at [UP9](https://up9.com/).</br>
Notice: The project has been renamed to Kubeshark and moved to https://github.com/kubeshark/kubeshark.</br>

---------------------
Podloxx was created by [mogenius](https://mogenius.com) - The Virtual DevOps platform
