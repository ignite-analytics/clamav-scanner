# ClamAV Scanner

[![Continuous Integration](https://github.com/ignite-analytics/clamav-scanner/actions/workflows/ci.yaml/badge.svg)](https://github.com/ignite-analytics/clamav-scanner/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ignite-analytics/clamav-scanner)](https://goreportcard.com/report/github.com/ignite-analytics/clamav-scanner)
[![Static Badge](https://img.shields.io/badge/licence%20-%20MIT-green)](https://github.com/ignite-analytics/clamav-scanner/blob/main/LICENSE)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/ignite-analytics/clamav-scanner)](https://github.com/ignite-analytics/clamav-scanner/releases)

*This project is largely inspired by the [Malware Scanner Service](https://github.com/GoogleCloudPlatform/docker-clamav-malware-scanner).*

It is a simple service that scans files for malware using the ClamAV antivirus engine. The difference is that this project is designed to be deployed on Google Cloud Run or GKE, and it is written in Go. It also comes with example Pulumi Typescript app and a Helm3 chart.

### How it works
The architecture is a bit different because there is no need to have Unscanned bucket used as sink for all uploaded files. Instead, you register buckets you whish to scan and if malware is found, the file is moved to Quarantine bucket. The image is build using Dockerfile and runs ClamAV as unprivilegued user and listens on port 1337.

**Google Kubernetes Engine**

![GKE diagram](gke-diagram.svg)

**Cloud Run**

![Cloud Run diagram](cloudrun-diagram.svg)

### Endpoints

There are 4 endpoints available:
- `/health` - returns 200 OK if `clamd` process is running. Used for health checks and Kubernetes probes.
- `/mirror` - serves gcp storage bucket as private mirror for database files. Kudos to [google-storage-proxy](https://github.com/cirruslabs/google-storage-proxy)
- `/update` - updates ClamAV database files. It accepts `GET` request and triggers and update. It can be used by Cloud Scheduler or Kubernetes CronJob.
- `/scan` - scans file for malware. It accepts `POST` request with following json payload:

```json
{
    "name": "object-name",
    "bucket": "bucket-name"
}
```

### Deployment

There are two ways to deploy this service and both examples are located in `infrastructure` directory:
- Helm3 chart for GKE deployment
- OpenTofu module for Cloud Run deployment
- Pulumi Typescript app for Cloud Run deployment

> [!IMPORTANT]
> Cloud Run supports only images from Artifact Registry and Docker Hub. ATM we only provide pre-build images as GitHub package which works fine for GKE. If you want to deploy on Cloud Run, you need to build the image and push it to Artifact Registry.

## Issues, Features, Feedback

Your input matters. Feel free to open [issues](https://github.com/ignite-analytics/clamav-scanner/issues) for bugs, feature requests, or any feedback you may have. Check if a similar issue exists before creating a new one, and please use clear titles and explanations to help understand your point better. Your thoughts help me improve this project!

### How to Contribute

🌟 Thank you for considering contributing to my project! Your efforts are incredibly valuable. To get started:

1. Fork the repository.
2. Create your feature branch: `git checkout -b feature/YourFeature`
3. Commit your changes: `git commit -am 'Add: YourFeature'`
4. Push to the branch: `git push origin feature/YourFeature`
5. Submit a pull request! 🚀
