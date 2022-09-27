<md- $bundle_name := .Name -md>
<md- $cloud_prefix := .CloudPrefix -md>
<md- $repo_name := .RepoName -md>
<md- $repo_encoded := .RepoNameEncoded -md>

[![Massdriver][logo]][website]

# <md $bundle_name md>

[![Release][release_shield]][release_url]
[![Contributors][contributors_shield]][contributors_url]
[![Forks][forks_shield]][forks_url]
[![Stargazers][stars_shield]][stars_url]
[![Issues][issues_shield]][issues_url]
[![MIT License][license_shield]][license_url]

<md .Description md>

---

## Design

For detailed information, check out our [Operator Guide](operator.mdx) for this bundle.

## Usage

Our bundles aren't intended to be used locally, outside of testing. Instead, our bundles are designed to be configured, connected, deployed and monitored in the [Massdriver][website] platform.

### What are Bundles?

Bundles are the basic building blocks of infrastructure, applications, and architectures in [Massdriver][website]. Read more [here](https://docs.massdriver.cloud/concepts/bundles).

## Bundle

<!-- COMPLIANCE:START -->

Security and compliance scanning of our bundles is performed using [Bridgecrew](https://www.bridgecrew.cloud/). Massdriver also offers security and compliance scanning of operational infrastructure configured and deployed using the platform.

| Benchmark                                                                                                                                                                                                                                                       | Description                        |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------- |
| [![Infrastructure Security](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/general" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=INFRASTRUCTURE+SECURITY" $repo_encoded md>) | Infrastructure Security Compliance |

<md if eq $cloud_prefix "k8s" -md>
| [![CIS KUBERNETES](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/cis_kubernetes" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=CIS+KUBERNETES+V1.5" $repo_encoded md>) | Center for Internet Security, KUBERNETES Compliance |
<md else if eq $cloud_prefix "aws" -md>
| [![CIS AWS](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/cis_aws" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=CIS+AWS+V1.2" $repo_encoded md>) | Center for Internet Security, AWS Compliance |
<md else if eq $cloud_prefix "azure" -md>
| [![CIS AZURE](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/cis_azure" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=CIS+AZURE+V1.1" $repo_encoded md>) | Center for Internet Security, AZURE Compliance |
<md else if eq $cloud_prefix "gcp" -md>
| [![CIS GCP](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/cis_gcp" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=CIS+GCP+V1.1" $repo_encoded md>) | Center for Internet Security, GCP Compliance |
<md end -md>
| [![PCI-DSS](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/pci" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=PCI-DSS+V3.2" $repo_encoded md>) | Payment Card Industry Data Security Standards Compliance |
| [![NIST-800-53](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/nist" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=NIST-800-53" $repo_encoded md>) | National Institute of Standards and Technology Compliance |
| [![ISO27001](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/iso" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=ISO27001" $repo_encoded md>) | Information Security Management System, ISO/IEC 27001 Compliance |
| [![SOC2](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/soc2" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=SOC2" $repo_encoded md>)| Service Organization Control 2 Compliance |
| [![HIPAA](<md printf "https://www.bridgecrew.cloud/badges/github/massdriver-cloud/%s/hipaa" $bundle_name md>)](<md printf "https://www.bridgecrew.cloud/link/badge?vcs=github&fullRepo=%s&benchmark=HIPAA" $repo_encoded md>) | Health Insurance Portability and Accountability Compliance |

<!-- COMPLIANCE:END -->

### Params

Form input parameters for configuring a bundle for deployment.

<details>
<summary>View</summary>

<!-- PARAMS:START -->

**Params coming soon**

<!-- PARAMS:END -->

</details>

### Connections

Connections from other bundles that this bundle depends on.

<details>
<summary>View</summary>

<!-- CONNECTIONS:START -->

**Connections coming soon**

<!-- CONNECTIONS:END -->

</details>

### Artifacts

Resources created by this bundle that can be connected to other bundles.

<details>
<summary>View</summary>

<!-- ARTIFACTS:START -->

**Artifacts coming soon**

<!-- ARTIFACTS:END -->

</details>

## Contributing

<!-- CONTRIBUTING:START -->

### Bug Reports & Feature Requests

Did we miss something? Please [submit an issue](<md printf "https://github.com/massdriver-cloud/%s/issues" $bundle_name md>) to report any bugs or request additional features.

### Developing

**Note**: Massdriver bundles are intended to be tightly use-case scoped, intention-based, reusable pieces of IaC for use in the [Massdriver][website] platform. For this reason, major feature additions that broaden the scope of an existing bundle are likely to be rejected by the community.

Still want to get involved? First check out our [contribution guidelines](https://docs.massdriver.cloud/bundles/contributing).

### Fix or Fork

If your use-case isn't covered by this bundle, you can still get involved! Massdriver is designed to be an extensible platform. Fork this bundle, or [create your own bundle from scratch](https://docs.massdriver.cloud/bundles/development)!

<!-- CONTRIBUTING:END -->

## Connect

<!-- CONNECT:START -->

Questions? Concerns? Adulations? We'd love to hear from you!

Please connect with us!

[![Email][email_shield]][email_url]
[![GitHub][github_shield]][github_url]
[![LinkedIn][linkedin_shield]][linkedin_url]
[![Twitter][twitter_shield]][twitter_url]
[![YouTube][youtube_shield]][youtube_url]
[![Reddit][reddit_shield]][reddit_url]

<md $utm_link := printf "%%s?utm_source=%s&utm_medium=%s&utm_campaign=%s&utm_content=%s" "github" "readme" $bundle_name "%s" -md>

<!-- markdownlint-disable -->

[logo]: https://raw.githubusercontent.com/massdriver-cloud/docs/main/static/img/logo-with-logotype-horizontal-400x110.svg

[docs]: <md printf $utm_link "https://docs.massdriver.cloud/" "docs" md>
[website]: <md printf $utm_link "https://www.massdriver.cloud/" "website" md>
[github]: <md printf $utm_link "https://github.com/massdriver-cloud" "github" md>
[slack]: <md printf $utm_link "https://massdriverworkspace.slack.com/" "slack" md>
[linkedin]: <md printf $utm_link "https://www.linkedin.com/company/massdriver/" "linkedin" md>

[contributors_shield]: <md printf "https://img.shields.io/github/contributors/massdriver-cloud/%s.svg?style=for-the-badge" $bundle_name md>
[contributors_url]: <md printf "https://github.com/massdriver-cloud/%s/graphs/contributors" $bundle_name md>
[forks_shield]: <md printf "https://img.shields.io/github/forks/massdriver-cloud/%s.svg?style=for-the-badge" $bundle_name md>
[forks_url]: <md printf "https://github.com/massdriver-cloud/%s/network/members" $bundle_name md>
[stars_shield]: <md printf "https://img.shields.io/github/stars/massdriver-cloud/%s.svg?style=for-the-badge" $bundle_name md>
[stars_url]: <md printf "https://github.com/massdriver-cloud/%s/stargazers" $bundle_name md>
[issues_shield]: <md printf "https://img.shields.io/github/issues/massdriver-cloud/%s.svg?style=for-the-badge" $bundle_name md>
[issues_url]: <md printf "https://github.com/massdriver-cloud/%s/issues" $bundle_name md>
[release_url]: <md printf "https://github.com/massdriver-cloud/%s/releases/latest" $bundle_name md>
[release_shield]: <md printf "https://img.shields.io/github/release/massdriver-cloud/%s.svg?style=for-the-badge" $bundle_name md>
[license_shield]: <md printf "https://img.shields.io/github/license/massdriver-cloud/%s.svg?style=for-the-badge" $bundle_name md>
[license_url]: <md printf "https://github.com/massdriver-cloud/%s/blob/main/LICENSE" $bundle_name md>

[email_url]: mailto:support@massdriver.cloud
[email_shield]: https://img.shields.io/badge/email-Massdriver-black.svg?style=for-the-badge&logo=mail.ru&color=000000
[github_url]: mailto:support@massdriver.cloud
[github_shield]: https://img.shields.io/badge/follow-Github-black.svg?style=for-the-badge&logo=github&color=181717
[linkedin_url]: https://linkedin.com/in/massdriver-cloud
[linkedin_shield]: https://img.shields.io/badge/follow-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&color=0A66C2

[twitter_url]: <md printf $utm_link "https://twitter.com/massdriver" "twitter" md>
[twitter_shield]: https://img.shields.io/badge/follow-Twitter-black.svg?style=for-the-badge&logo=twitter&color=1DA1F2
[discourse_url]: <md printf $utm_link "https://community.massdriver.cloud" "discourse" md>
[discourse_shield]: https://img.shields.io/badge/join-Discourse-black.svg?style=for-the-badge&logo=discourse&color=000000
[youtube_url]: https://www.youtube.com/channel/UCfj8P7MJcdlem2DJpvymtaQ
[youtube_shield]: https://img.shields.io/badge/subscribe-Youtube-black.svg?style=for-the-badge&logo=youtube&color=FF0000
[reddit_url]: https://www.reddit.com/r/massdriver
[reddit_shield]: https://img.shields.io/badge/subscribe-Reddit-black.svg?style=for-the-badge&logo=reddit&color=FF4500

<!-- markdownlint-restore -->

<!-- CONNECT:END -->
