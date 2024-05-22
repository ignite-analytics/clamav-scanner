# ClamAV Scanner

*This project is largely inspired by the [Malware Scanner Service](https://github.com/GoogleCloudPlatform/docker-clamav-malware-scanner).*

It is a simple service that scans files for malware using the ClamAV antivirus engine. The difference is that this project is designed to be deployed on Google Cloud Run or GKE, and it is written in Go. It also comes with example Pulumi Typescript app and a Helm3 chart.

### How it works
The architecture is a bit different because there is no need to have Unscanned bucket used as sink for all uploaded files. Instead, you register buckets you whish to scan and if malware is found, the file is moved to Quarantine bucket. The image is build using Dockerfile and runs ClamAV as unprivilegued user and listens on port 1337.

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
