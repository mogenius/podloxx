<p align="center">
<img src="ui/src/assets/logos/logo-horizontal.svg" alt="drawing" width="500"/>
    </p>
    
# podloxx: Skinny Kubernetes Traffic Viewer<br /> (including local traffic)

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

## What is podloxx?
podloxx is a easy to use cluster daemonset and frontend application to monitor the ongoing traffic of your cluster in real-time. The monitoring happens in real-time and the display happens every few seconds (for perfromance reasons).
The most important goal was to tell local and external traffic appart. This allows you to identify high traffic applications and see how is communicating with the outside world a lot. This even works with slim containers :-)

## Installation
Just download it and run it. Don't forget to set the right cluster using kubectx or whatever tool you prefer.

### Download
```
curl -Lo podloxx \
https://github.com/mogenius/podlox/releases/download/s/podloxx-1.0.2-darwin-arm64 \
&& chmod 755 podloxx

podloxx start
```

## How does it work? What does it do?
TODO Explain architecture

## Configuration
INTERFACE_PREFIX, azv, veth, cali

## API
TODO

## Roadmap
TODO 

## Known Issues
TODO

## Thanks 
We took great inspiration (and some lines of code) from [Mizu](https://github.com/up9inc/mizu). Awesome work from the folks at [UP9](https://up9.com/).

## About mogenius
mogenius provides an automated cloud infrastructure that allows scaling applications on Kubernetes with a user friendly UI and API. 
