<div id="top"></div>
<div align="center">
  <h1 align="center">Vaxiin CLI</h1>

  <p align="center">
    CLI component of the Vaxiin framework
    <br />
    <a href="https://docs.vaxiin.io"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/rebootoio/vaxctl/issues/new?assignees=&labels=bug&template=bug_report.md&title=">Report Bug</a>
    ·
    <a href="https://github.com/rebootoio/vaxctl/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=">Request Feature</a>
  </p>
    <a >
    <img src='https://img.shields.io/github/v/tag/rebootoio/vaxctl?style=for-the-badge'>
  </a>
  <a href='https://discord.gg/aEJ6qwcCGs'>
    <img src='https://img.shields.io/discord/813371439469297674?style=for-the-badge'>
  </a>
  <a href='https://github.com/rebootoio/vaxctl/blob/main/LICENSE'>
    <img src='https://img.shields.io/github/license/rebootoio/vaxctl?style=for-the-badge'>
  </a>
</div>

## Quick Start Guide
To get started quickly with all the component, follow the guide at the [sandbox repository](https://github.com/rebootoio/vaxiin-sandbox)


## Components
| Component | Repoistory | Artifact | Documentation |
|-----------|------------|----------|------|
| Server | [GitHub](https://github.com/rebootoio/vaxiin-server) | [DockerHub](https://hub.docker.com/repository/docker/rebooto/vaxiin-server) | [Docs](https://docs.vaxiin.io/configuration/server) |
| Handler | [GitHub](https://github.com/rebootoio/vaxiin-handler) | [DockerHub](https://hub.docker.com/repository/docker/rebooto/vaxiin-handler) | [Docs](https://docs.vaxiin.io/configuration/handler) |
| Agent | [GitHub](https://github.com/rebootoio/vaxiin-agent) | [DockerHub](https://hub.docker.com/repository/docker/rebooto/vaxiin-agent) | [Docs](https://docs.vaxiin.io/configuration/agent) |
| CLI | [GitHub](https://github.com/rebootoio/vaxctl)| [Release](https://github.com/rebootoio/vaxctl/releases) | [Docs](https://docs.vaxiin.io/configuration/cli) |

## Running
Go to the [releases page](https://github.com/rebootoio/vaxctl/releases), download the latest binary and place it in a directory you have in your `PATH`

Download the sample configuration file locally
```
wget https://raw.githubusercontent.com/rebootoio/vaxctl/main/config-example.yaml -O $HOME/.vaxctl.yaml
```
Edit the file and change the `url` so that it will  point to your vaxiin server

Start using the cli
```
vaxctl -h
```
_Additional information about the available commands can be found in the [docs](https://docs.vaxiin.io/cli-reference/overview)._

## Builiding
The go binary can be built by cloning this repository and running `make build`

## Contributing

Contributions are what make the Open Source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**!

If you have a suggestion that would make this better, please fork the repo and create a Pull Request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


## License
Distributed under the [AGPL-3.0 License](https://github.com/rebootoio/vaxctl/blob/main/LICENSE) License.

## Contact
[Join the Rebooto Discord Server](https://discord.gg/aEJ6qwcCGs)

[Open an issue](https://github.com/rebootoio/vaxctl/issues)
