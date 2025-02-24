Here you can find example of Pulumi Typescript application to deploy the service on Google Cloud Run. All you need to do
is to configure your variables in file `globalConf.ts` and run `pulumi up`.

### What is created

- Service Account and adds IAM permissions for Eventarc events
- Creates mirror and quarantine buckets. Bucket used for logging is expect to be already created
- Creates KMS keyring and crypto key used for encrypting mirror and quarantine buckets
- Creates Cloud Run service 2nd generation
- Creates Cloud Scheduler job that calls `update` endpoint
- Creates Eventarc trigger that listens for `google.cloud.storage.object.v1.finalized` events. This is created by
  populating `event.buckets` in `globalConf.ts`. It's a list of objects where you declare name and location. That is the
  place where you register buckets to be scanned for malware. Event triggers are created in same location as the
  buckets.

### Testing

This application uses [`gcp-pac`](https://github.com/losisin/gcp-pac) for Policy-as-Code and all policies are set to
mandatory. It also comes with basic set of unit tests to get you started.
