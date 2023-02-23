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
