export const name = 'clamav-scanner' // service name
export const project = 'my-project' // GCP project ID
export const region = 'europe-west1' // service location
export const image = {
	repository: 'ghcr.io/ignite-analytics/clamav-scanner',
	tag: '0.3.2'
}
export const schedule = {
	cron: '37 */2 * * *',
	timezone: 'Europe/Oslo',
	location: 'europe-west1' // cheduler job location
}
export const service = {
	port: 1337
}
export const event = {
	type: 'google.cloud.storage.object.v1.finalized',
	buckets: [
		{
			name: 'bucket-name',
			location: 'EU'
		}
	]
}
export const storage = {
	location: 'EU',
	mirrorBucket: `${name}-mirror`,
	quarantineBucket: `${name}-quarantine`,
	logBucket: 'my-log-bucket' // where to store bucket logs
}
export const kms = {
	location: 'europe', // KMS keyring location
	rotationPeriod: '7776000s' // 90 days
}
export const labels = {
	manager: 'pulumi',
	service: name
}
export const pubsub = {
	topic: `${name}-topic`,
	subscription: `${name}-subscription`,
	subscribers: ['<SERVICE_NAME>@<YOUR_PROJECT_ID>.iam.gserviceaccount.com'],
	gcsServiceAccount: 'service-<YOUR_PROJECT_NUMBER>@gcp-sa-pubsub.iam.gserviceaccount.com'
}
